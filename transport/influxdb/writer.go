package influxdb

import (
	"bytes"
	"context"
	"fmt"
	"time"

	influx "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/rs/zerolog/log"
	"github.com/vodolaz095/dashboard/model"
	"github.com/vodolaz095/dashboard/service"
)

type Writer struct {
	Endpoint     string
	Token        string
	Organization string
	Bucket       string
	Service      *service.SensorsService

	initialized bool
	client      influx.Client
	writer      api.WriteAPI
}

func (w *Writer) Init(ctx context.Context) (err error) {
	w.client = influx.NewClient(w.Endpoint, w.Token)
	w.writer = w.client.WriteAPI(w.Organization, w.Bucket)
	w.initialized = true
	_, err = w.client.Ping(ctx)
	return
}

func (w *Writer) Ping(ctx context.Context) (err error) {
	if !w.initialized {
		return nil
	}
	_, err = w.client.Ping(ctx)
	return
}

func (w *Writer) format(input model.Update) string {
	// stat,unit=temperature,a=b value=%f 1556813561098000000
	if input.Error != "" {
		return ""
	}
	sensor, found := w.Service.Sensors[input.Name]
	if !found {
		return ""
	}
	ret := bytes.NewBufferString(sensor.GetName())
	for k, v := range sensor.GetTags() {
		fmt.Fprintf(ret, ",%s=%s", k, v)
	}
	fmt.Fprintf(ret, ",type=%s", sensor.GetType())
	fmt.Fprintf(ret, " value=%f %v\n", input.Value, time.Now().UnixNano())
	return ret.String()
}

func (w *Writer) Start(ctx context.Context) {
	if !w.initialized {
		return
	}
	feed, err := w.Service.Subscribe(ctx, "dashboard.broadcaster.influxdb")
	if err != nil {
		log.Fatal().Err(err).Msgf("broadcaster failed to subscribe: %s", err)
	}
	for {
		select {
		case <-ctx.Done():
			log.Info().Msgf("Influxdb writer is closing...")
			w.writer.Flush()
			w.client.Close()
			return

		case upd := <-feed:
			w.writer.WriteRecord(w.format(upd))
		}
	}
}
