package main

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	redisClient "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/vodolaz095/dashboard/transport/broadcaster"
	"github.com/vodolaz095/dqueue"

	"github.com/vodolaz095/dashboard/config"
	"github.com/vodolaz095/dashboard/pkg/healthcheck"
	"github.com/vodolaz095/dashboard/pkg/zerologger"
	"github.com/vodolaz095/dashboard/sensors"
	"github.com/vodolaz095/dashboard/sensors/curl"
	"github.com/vodolaz095/dashboard/sensors/endpoint"
	"github.com/vodolaz095/dashboard/sensors/mysql"
	"github.com/vodolaz095/dashboard/sensors/postgres"
	"github.com/vodolaz095/dashboard/sensors/redis"
	"github.com/vodolaz095/dashboard/sensors/shell"
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
	var connectionIsFound bool
	byIndex := make([]string, len(cfg.Sensors))
	byName := make(map[string]sensors.ISensor, 0)
	for i := range cfg.Sensors {
		log.Debug().
			Msgf("Setting up sensor %v: %s of type %s", i,
				cfg.Sensors[i].Name, cfg.Sensors[i].Type,
			)
		switch cfg.Sensors[i].Type {
		case "mysql", "mariadb":
			ms := &mysql.Sensor{}
			ms.Name = cfg.Sensors[i].Name
			ms.Type = "mysql"
			ms.RefreshRate = cfg.Sensors[i].RefreshRate
			ms.Description = cfg.Sensors[i].Description
			ms.Link = cfg.Sensors[i].Link
			ms.Minimum = cfg.Sensors[i].Minimum
			ms.Maximum = cfg.Sensors[i].Maximum
			ms.Tags = cfg.Sensors[i].Tags

			ms.Query = cfg.Sensors[i].Query
			ms.DatabaseConnectionName = cfg.Sensors[i].ConnectionName
			ms.Con, connectionIsFound = srv.MysqlConnections[cfg.Sensors[i].ConnectionName]
			if !connectionIsFound {
				log.Fatal().Msgf("Sensor %v %s refers to unknown %s connection %s",
					i, cfg.Sensors[i].Name, cfg.Sensors[i].Type, cfg.Sensors[i].ConnectionName,
				)
			}
			byName[ms.Name] = ms
			byIndex[i] = ms.Name
			updateQueue.ExecuteAfter(ms.Name, 50*time.Millisecond)
			break
		case "redis":
			rs := &redis.Sensor{}
			rs.Name = cfg.Sensors[i].Name
			rs.Type = "redis"
			rs.RefreshRate = cfg.Sensors[i].RefreshRate
			rs.Description = cfg.Sensors[i].Description
			rs.Link = cfg.Sensors[i].Link
			rs.Minimum = cfg.Sensors[i].Minimum
			rs.Maximum = cfg.Sensors[i].Maximum
			rs.Tags = cfg.Sensors[i].Tags

			rs.Query = cfg.Sensors[i].Query
			rs.DatabaseConnectionName = cfg.Sensors[i].ConnectionName
			rs.Client, connectionIsFound = srv.RedisConnections[cfg.Sensors[i].ConnectionName]
			if !connectionIsFound {
				log.Fatal().Msgf("Sensor %v %s refers to unknown %s connection %s",
					i, cfg.Sensors[i].Name, cfg.Sensors[i].Type, cfg.Sensors[i].ConnectionName,
				)
			}

			byName[rs.Name] = rs
			byIndex[i] = rs.Name
			updateQueue.ExecuteAfter(rs.Name, 50*time.Millisecond)

			break
		case "postgres":
			ps := &postgres.Sensor{}
			ps.Name = cfg.Sensors[i].Name
			ps.Type = "postgres"
			ps.RefreshRate = cfg.Sensors[i].RefreshRate
			ps.Description = cfg.Sensors[i].Description
			ps.Link = cfg.Sensors[i].Link
			ps.Minimum = cfg.Sensors[i].Minimum
			ps.Maximum = cfg.Sensors[i].Maximum
			ps.Tags = cfg.Sensors[i].Tags

			ps.Query = cfg.Sensors[i].Query
			ps.DatabaseConnectionName = cfg.Sensors[i].ConnectionName
			ps.Con, connectionIsFound = srv.PostgresqlConnections[cfg.Sensors[i].ConnectionName]
			if !connectionIsFound {
				log.Fatal().Msgf("Sensor %v %s refers to unknown %s connection %s",
					i, cfg.Sensors[i].Name, cfg.Sensors[i].Type, cfg.Sensors[i].ConnectionName,
				)
			}

			byName[ps.Name] = ps
			byIndex[i] = ps.Name
			updateQueue.ExecuteAfter(ps.Name, 50*time.Millisecond)

			break
		case "curl":
			cs := &curl.Sensor{}
			cs.Name = cfg.Sensors[i].Name
			cs.Type = "curl"
			cs.RefreshRate = cfg.Sensors[i].RefreshRate
			cs.Description = cfg.Sensors[i].Description
			cs.Link = cfg.Sensors[i].Link
			cs.Minimum = cfg.Sensors[i].Minimum
			cs.Maximum = cfg.Sensors[i].Maximum
			cs.Tags = cfg.Sensors[i].Tags

			cs.HttpMethod = cfg.Sensors[i].HttpMethod
			cs.Endpoint = cfg.Sensors[i].Endpoint
			cs.Headers = cfg.Sensors[i].Headers
			cs.Body = cfg.Sensors[i].Body
			cs.JsonPath = cfg.Sensors[i].JsonPath

			byName[cs.Name] = cs
			byIndex[i] = cs.Name
			updateQueue.ExecuteAfter(cs.Name, 50*time.Millisecond)

			break
		case "shell":
			ss := &shell.Sensor{}
			ss.Name = cfg.Sensors[i].Name
			ss.Type = "shell"
			ss.RefreshRate = cfg.Sensors[i].RefreshRate
			ss.Description = cfg.Sensors[i].Description
			ss.Link = cfg.Sensors[i].Link
			ss.Minimum = cfg.Sensors[i].Minimum
			ss.Maximum = cfg.Sensors[i].Maximum
			ss.Tags = cfg.Sensors[i].Tags

			ss.Command = cfg.Sensors[i].Command
			ss.Environment = cfg.Sensors[i].Environment
			ss.JsonPath = cfg.Sensors[i].JsonPath

			byName[ss.Name] = ss
			byIndex[i] = ss.Name
			updateQueue.ExecuteAfter(ss.Name, 50*time.Millisecond)

			break
		case "endpoint":
			es := &endpoint.Sensor{}
			es.Name = cfg.Sensors[i].Name
			es.Type = "endpoint"
			es.RefreshRate = cfg.Sensors[i].RefreshRate
			es.Description = cfg.Sensors[i].Description
			es.Link = cfg.Sensors[i].Link
			es.Minimum = cfg.Sensors[i].Minimum
			es.Maximum = cfg.Sensors[i].Maximum
			es.Tags = cfg.Sensors[i].Tags
			es.Token = cfg.Sensors[i].Token

			byName[es.Name] = es
			byIndex[i] = es.Name
			updateQueue.ExecuteAfter(es.Name, 50*time.Millisecond)

			break
		default:
			log.Fatal().Msgf("Element %v has unknown sensor type %s", i, cfg.Sensors[i].Type)
		}
	}

	srv.ListOfSensors = byIndex
	srv.Sensors = byName
	// init service
	err = srv.InitSensors(ctx)
	if err != nil {
		log.Fatal().Err(err).Msgf("error initializing sensors: %s", err)
	}
	log.Info().Msgf("Sensors (%v) initialized", len(srv.Sensors))

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
		err = healthcheck.SetStatus("Dashboard is online!")
		if err != nil {
			log.Warn().Err(err).Msgf("Error setting systemd unit status")
		}
	}

	// configure redis broadcaster
	publisher := broadcaster.Publisher{
		Service: &srv,
	}
	if len(cfg.Broadcasters) > 0 {
		for i := range cfg.Broadcasters {
			err = publisher.InitConnection(cfg.Broadcasters[i].ConnectionName,
				cfg.Broadcasters[i].Subject)
			if err != nil {
				log.Fatal().Err(err).Msgf("Error initializing broadcaster for connection %v - %s",
					i, err,
				)
			}
		}
	}

	// configure webserver transport
	webServerTransport := webserver.Transport{
		Address:     cfg.Listen,
		Version:     Version,
		Domain:      cfg.Domain,
		Title:       cfg.Title,
		Description: cfg.Description,
		Keywords:    cfg.Keywords,
		DoIndex:     cfg.DoIndex,

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
	go publisher.Start(ctx)

	go func() {
		log.Debug().Msgf("Preparing to start webserver on %s...", webServerTransport.Address)
		errWeb := webServerTransport.Start(ctx, &wg)
		if errWeb != nil {
			log.Fatal().Err(errWeb).
				Msgf("error starting webserver on %s : %s", webServerTransport.Address, errWeb)
		}
	}()

	// todo: redis subscribber
	// todo: mqtt subscribber

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
