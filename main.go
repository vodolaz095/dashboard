package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	redisClient "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/vodolaz095/dashboard/transport/influxdb"
	redis_transport "github.com/vodolaz095/dashboard/transport/redis"
	"github.com/vodolaz095/dqueue"

	"github.com/vodolaz095/dashboard/config"
	"github.com/vodolaz095/dashboard/pkg/healthcheck"
	"github.com/vodolaz095/dashboard/pkg/zerologger"
	"github.com/vodolaz095/dashboard/sensors"
	"github.com/vodolaz095/dashboard/service"
	"github.com/vodolaz095/dashboard/transport/webserver"
)

var Version = "development"

func main() {
	var err error
	wg := sync.WaitGroup{}
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	flag.Parse()

	// load config
	if len(flag.Args()) != 1 {
		log.Fatal().Msgf("please, provide path to config as 1st argument")
	}
	pathToConfig := flag.Args()[0]
	cfg, err := config.LoadFromFile(pathToConfig)
	if err != nil {
		log.Fatal().Err(err).
			Msgf("error loading config from %s: %s", pathToConfig, err)
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(cfg)
	if err != nil {
		log.Fatal().Err(err).
			Msgf("error validating configuration file %s: %s", pathToConfig, err)
	}

	// set logging
	zerologger.Configure(cfg.Log)
	log.Debug().Msgf("Configuring %v sensors...", len(cfg.Sensors))

	if len(cfg.Sensors) == 0 {
		log.Fatal().Msgf("No sensors configured in %s!", pathToConfig)
	}

	// init service
	updateQueue := dqueue.New()
	srv := service.SensorsService{
		UpdateQueue:           &updateQueue,
		UpdateInterval:        100 * time.Millisecond,
		MysqlConnections:      make(map[string]*sql.Conn, 0),
		PostgresqlConnections: make(map[string]*sql.Conn, 0),
		RedisConnections:      make(map[string]*redisClient.Client, 0),
	}
	// init database connections
	for i := range cfg.DatabaseConnections {
		err = srv.MakeConnection(ctx,
			cfg.DatabaseConnections[i].Name,
			service.DatabaseConnectionType(cfg.DatabaseConnections[i].Type),
			cfg.DatabaseConnections[i].DatabaseConnectionString,
		)
		if err != nil {
			log.Fatal().Err(err).Msgf("Error initializing connection %s (%s): %s",
				cfg.DatabaseConnections[i].Name,
				cfg.DatabaseConnections[i].Type,
				err,
			)
		}
	}
	log.Info().Msgf("Database connections (%v) initialized", len(cfg.DatabaseConnections))

	// generate sensors from config
	byIndex := make([]string, len(cfg.Sensors))
	byName := make(map[string]sensors.ISensor, 0)
	var sensorCreated sensors.ISensor
	for i := range cfg.Sensors {
		log.Debug().
			Msgf("Setting up sensor %v: %s of type %s", i,
				cfg.Sensors[i].Name, cfg.Sensors[i].Type,
			)
		_, duplicateSensorFound := byName[cfg.Sensors[i].Name]
		if duplicateSensorFound {
			log.Fatal().Msgf("Sensor with duplicate name %s is found at index %v in config file",
				cfg.Sensors[i].Name, i,
			)
		}
		sensorCreated, err = srv.MakeSensor(cfg.Sensors[i])
		if err != nil {
			log.Fatal().Err(err).
				Msgf("Error creating sensor %v %s of type %s : %s",
					i, cfg.Sensors[i].Name, cfg.Sensors[i].Type, err)
		}
		byName[sensorCreated.GetName()] = sensorCreated
		byIndex[i] = sensorCreated.GetName()
	}

	srv.ListOfSensors = byIndex
	srv.Sensors = byName
	// init service
	err = srv.InitSensors(ctx)
	if err != nil {
		log.Fatal().Err(err).Msgf("error initializing sensors: %s", err)
	}
	log.Info().Msgf("Sensors (%v) initialized", len(srv.Sensors))

	// configure influxdb transport to store persistent timeseries data
	influxWriter := &influxdb.Writer{
		Endpoint:     cfg.Influx.Endpoint,
		Token:        cfg.Influx.Token,
		Organization: cfg.Influx.Organization,
		Bucket:       cfg.Influx.Bucket,
		Service:      &srv,
	}
	if cfg.Influx.Valid() {
		err = influxWriter.Init(ctx)
		if err != nil {
			log.Fatal().Err(err).
				Msgf("Error initializing influxdb connection to %s (org=%s, bucket=%s): %s",
					err, cfg.Influx.Endpoint, cfg.Influx.Organization, cfg.Influx.Bucket,
				)
		}
	}

	// set systemd watchdog
	systemdWatchdogEnabled, err := healthcheck.Ready()
	if err != nil {
		log.Error().Err(err).
			Msgf("%s: while notifying systemd on application ready", err)
	}
	if systemdWatchdogEnabled {
		go func() {
			log.Debug().Msgf("Watchdog enabled")
			errWD := healthcheck.StartWatchDog(ctx, []healthcheck.Pinger{
				&srv,
				influxWriter,
			})
			if errWD != nil {
				log.Error().
					Err(err).
					Msgf("%s : while starting watchdog", err)
			}
		}()
	} else {
		log.Warn().Msgf("Systemd watchdog disabled - application can work unstable in systemd environment")
	}

	// change systemd status
	if systemdWatchdogEnabled {
		// https://www.freedesktop.org/software/systemd/man/latest/sd_notify.html#STATUS=%E2%80%A6
		err = healthcheck.SetStatus(fmt.Sprintf(
			"Dashboard is online - %v sensors are monitored!", len(srv.Sensors)),
		)
		if err != nil {
			log.Warn().Err(err).Msgf("Error setting systemd unit status")
		}
	}

	// configure redis transports
	redisPublisher := redis_transport.Publisher{
		Service: &srv,
	}
	if len(cfg.Broadcasters) > 0 {
		for i := range cfg.Broadcasters {
			err = redisPublisher.InitConnection(
				cfg.Broadcasters[i].ConnectionName,
				cfg.Broadcasters[i].Subject,
				cfg.Broadcasters[i].ValueOnly,
			)
			if err != nil {
				log.Fatal().Err(err).Msgf("Error initializing broadcaster for connection %v - %s",
					i, err,
				)
			}
		}
	}
	redisSubscriber := redis_transport.Subscriber{
		Service: &srv,
	}

	// configure webserver transport
	webServerTransport := webserver.Transport{
		Address:        cfg.WebUI.Listen,
		Version:        Version,
		Domain:         cfg.WebUI.Domain,
		Title:          cfg.WebUI.Title,
		Description:    cfg.WebUI.Description,
		Keywords:       cfg.WebUI.Keywords,
		DoIndex:        cfg.WebUI.DoIndex,
		PathToHeader:   cfg.WebUI.PathToHeader,
		PathToFooter:   cfg.WebUI.PathToFooter,
		SensorsService: &srv,
	}

	// handle signals
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGABRT,
	)
	go func() {
		s := <-sigc
		log.Info().Msgf("Signal %s is received", s.String())
		wg.Done()
		cancel()
	}()

	// main loop
	go srv.StartRefreshingSensors(ctx)
	go srv.StartClock(ctx)
	go redisPublisher.Start(ctx)
	go redisSubscriber.Start(ctx)
	go influxWriter.Start(ctx)

	go func() {
		log.Debug().Msgf("Preparing to start webserver on %s...", webServerTransport.Address)
		errWeb := webServerTransport.Start(ctx, &wg)
		if errWeb != nil {
			log.Fatal().Err(errWeb).
				Msgf("error starting webserver on %s : %s", webServerTransport.Address, errWeb)
		}
	}()

	// todo: mqtt subscriber

	wg.Wait()
	terminationContext, terminationContextCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer terminationContextCancel()
	err = srv.Close(terminationContext)
	if err != nil {
		log.Error().Err(err).
			Msgf("Error terminating application, something can be broken: %s", err)
	}
	log.Info().Msgf("Dashboard is terminated.")
}
