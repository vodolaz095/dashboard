package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"github.com/vodolaz095/dashboard/pkg/healthcheck"
	"github.com/vodolaz095/dashboard/pkg/zerologger"
	"github.com/vodolaz095/dashboard/sensors"
	"github.com/vodolaz095/dashboard/transport/webserver"
	"github.com/vodolaz095/dqueue"

	"github.com/vodolaz095/dashboard/config"
	"github.com/vodolaz095/dashboard/sensors/curl"
	"github.com/vodolaz095/dashboard/sensors/mysql"
	"github.com/vodolaz095/dashboard/sensors/postgres"
	"github.com/vodolaz095/dashboard/sensors/redis"
	"github.com/vodolaz095/dashboard/sensors/shell"
	"github.com/vodolaz095/dashboard/service"
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

	// generate sensors from config
	updateQueue := dqueue.New()
	byIndex := make([]string, len(cfg.Sensors))
	byName := make(map[string]sensors.ISensor, 0)
	for i := range cfg.Sensors {
		log.Debug().
			Msgf("Setting up sensor %v: %s of type %s", i,
				cfg.Sensors[i].Name, cfg.Sensors[i].Type,
			)
		switch cfg.Sensors[i].Type {
		case "mysql":
			ms := &mysql.Sensor{}
			ms.Name = cfg.Sensors[i].Name
			ms.Type = "mysql"
			ms.RefreshRate = cfg.Sensors[i].RefreshRate
			ms.Description = cfg.Sensors[i].Description
			ms.Link = cfg.Sensors[i].Link
			ms.Minimum = cfg.Sensors[i].Minimum
			ms.Maximum = cfg.Sensors[i].Maximum
			ms.Tags = cfg.Sensors[i].Tags

			ms.DatabaseConnectionString = cfg.Sensors[i].DatabaseConnectionString
			ms.Query = cfg.Sensors[i].Query

			byName[ms.Name] = ms
			byIndex[i] = ms.Name
			updateQueue.ExecuteAfter(ms.Name, 50*time.Millisecond)
			break
		case "redis":
			ms := &redis.Sensor{}
			ms.Name = cfg.Sensors[i].Name
			ms.Type = "redis"
			ms.RefreshRate = cfg.Sensors[i].RefreshRate
			ms.Description = cfg.Sensors[i].Description
			ms.Link = cfg.Sensors[i].Link
			ms.Minimum = cfg.Sensors[i].Minimum
			ms.Maximum = cfg.Sensors[i].Maximum
			ms.Tags = cfg.Sensors[i].Tags

			ms.DatabaseConnectionString = cfg.Sensors[i].DatabaseConnectionString
			ms.Query = cfg.Sensors[i].Query

			byName[ms.Name] = ms
			byIndex[i] = ms.Name
			updateQueue.ExecuteAfter(ms.Name, 50*time.Millisecond)

			break
		case "postgres":
			ms := &postgres.Sensor{}
			ms.Name = cfg.Sensors[i].Name
			ms.Type = "postgres"
			ms.RefreshRate = cfg.Sensors[i].RefreshRate
			ms.Description = cfg.Sensors[i].Description
			ms.Link = cfg.Sensors[i].Link
			ms.Minimum = cfg.Sensors[i].Minimum
			ms.Maximum = cfg.Sensors[i].Maximum
			ms.Tags = cfg.Sensors[i].Tags

			ms.DatabaseConnectionString = cfg.Sensors[i].DatabaseConnectionString
			ms.Query = cfg.Sensors[i].Query

			byName[ms.Name] = ms
			byIndex[i] = ms.Name
			updateQueue.ExecuteAfter(ms.Name, 50*time.Millisecond)

			break
		case "curl":
			ms := &curl.Sensor{}
			ms.Name = cfg.Sensors[i].Name
			ms.Type = "curl"
			ms.RefreshRate = cfg.Sensors[i].RefreshRate
			ms.Description = cfg.Sensors[i].Description
			ms.Link = cfg.Sensors[i].Link
			ms.Minimum = cfg.Sensors[i].Minimum
			ms.Maximum = cfg.Sensors[i].Maximum
			ms.Tags = cfg.Sensors[i].Tags

			ms.HttpMethod = cfg.Sensors[i].HttpMethod
			ms.Endpoint = cfg.Sensors[i].Endpoint
			ms.Headers = cfg.Sensors[i].Headers
			ms.Body = cfg.Sensors[i].Body
			ms.JsonPath = cfg.Sensors[i].JsonPath

			byName[ms.Name] = ms
			byIndex[i] = ms.Name
			updateQueue.ExecuteAfter(ms.Name, 50*time.Millisecond)

			break
		case "shell":
			ms := &shell.Sensor{}
			ms.Name = cfg.Sensors[i].Name
			ms.Type = "shell"
			ms.RefreshRate = cfg.Sensors[i].RefreshRate
			ms.Description = cfg.Sensors[i].Description
			ms.Link = cfg.Sensors[i].Link
			ms.Minimum = cfg.Sensors[i].Minimum
			ms.Maximum = cfg.Sensors[i].Maximum
			ms.Tags = cfg.Sensors[i].Tags

			ms.Command = cfg.Sensors[i].Command
			ms.Environment = cfg.Sensors[i].Environment
			ms.JsonPath = cfg.Sensors[i].JsonPath

			byName[ms.Name] = ms
			byIndex[i] = ms.Name
			updateQueue.ExecuteAfter(ms.Name, 50*time.Millisecond)

			break
		case "endpoint":
			ms := &shell.Sensor{}
			ms.Name = cfg.Sensors[i].Name
			ms.Type = "endpoint"
			ms.RefreshRate = cfg.Sensors[i].RefreshRate
			ms.Description = cfg.Sensors[i].Description
			ms.Link = cfg.Sensors[i].Link
			ms.Minimum = cfg.Sensors[i].Minimum
			ms.Maximum = cfg.Sensors[i].Maximum
			ms.Tags = cfg.Sensors[i].Tags

			ms.Token = cfg.Sensors[i].Token

			byName[ms.Name] = ms
			byIndex[i] = ms.Name
			updateQueue.ExecuteAfter(ms.Name, 50*time.Millisecond)

			break
		default:
			log.Fatal().Msgf("Element %v has unknown sensor type %s", i, cfg.Sensors[i].Type)
		}
	}

	// init service
	srv := service.SensorsService{
		UpdateQueue:    &updateQueue,
		UpdateInterval: 100 * time.Millisecond,
		ListOfSensors:  byIndex,
		Sensors:        byName,
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

	// configure transports
	webServerTransport := webserver.Transport{
		Address:        cfg.Listen,
		SensorsService: &srv,
		Domain:         cfg.Domain,
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
	go srv.StartKeepingSensorsUpToDate(ctx)

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
