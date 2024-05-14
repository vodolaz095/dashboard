package healthcheck

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/go-systemd/v22/daemon"
	"github.com/rs/zerolog/log"
)

// Pinger is interface StartWatchDog uses
type Pinger interface {
	Ping(context.Context) error
}

// Ready is used to nofity system is ready and check, if we need to perform systemd healthchecks
func Ready() (supported bool, err error) {
	return daemon.SdNotify(false, daemon.SdNotifyReady)
}

// SetStatus sets systemd unit status in human readable format
// https://www.freedesktop.org/software/systemd/man/latest/sd_notify.html#STATUS=%E2%80%A6
func SetStatus(status string) (err error) {
	_, err = daemon.SdNotify(false, fmt.Sprintf("STATUS=%s", status))
	return
}

// StartWatchDog starts background process that notifies systemd if application is running properly
func StartWatchDog(mainCtx context.Context, pingers []Pinger) (err error) {
	var ok bool
	interval, err := daemon.SdWatchdogEnabled(false)
	if err != nil {
		return
	}
	if interval == 0 {
		log.Info().Msgf("Watchdog not enabled")
		return
	}
	ticker := time.NewTicker(interval / 2)
	go func() {
		<-mainCtx.Done()
		ticker.Stop()
	}()
	for t := range ticker.C {
		ctx, cancel := context.WithDeadline(mainCtx, t.Add(interval/2))
		ok = true
		for i := range pingers {
			err = pingers[i].Ping(ctx)
			if err != nil {
				log.Error().Err(err).Msgf("%s: while pinging", err)
				ok = false
			}
		}
		cancel()
		if ok {
			_, err = daemon.SdNotify(false, daemon.SdNotifyWatchdog)
			if err != nil {
				log.Error().Err(err).Msgf("%s: while sending watchdog notification", err)
			}
			log.Trace().Msgf("Service is healthy!")
		} else {
			log.Warn().Msgf("Service is broken!")
		}
	}
	return nil
}
