package victoriametrics

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/vodolaz095/dashboard/internal/sensors"
)

// VMSenor executes periodical PromQL queries on remote Victoria Metrics Database
type VMSenor struct {
	sensors.UnimplementedSensor
	// Client is http.Client used to dial Victoria <etrics api
	Client *http.Client
	// Endpoint is used to dial database
	Endpoint string
	// Query is send to remote database in order to receive data from it
	Query string
	// Headers are HTTP request headers being send with any HTTP request to Victoria Metrics
	Headers map[string]string
	// Filter is used to select timeseries required
	Filter map[string]string
}

func (s *VMSenor) Init(ctx context.Context) error {
	s.Mutex = &sync.RWMutex{}
	if s.A == 0 {
		s.A = 1
	}
	s.Client = http.DefaultClient
	return s.Ping(ctx)
}

// Ping defines means to ensure victoria metrics is running
func (s *VMSenor) Ping(ctx context.Context) error {
	// https://github.com/VictoriaMetrics/VictoriaMetrics/issues/3539#issuecomment-1366469760
	req, err := http.NewRequest(http.MethodGet, s.Endpoint+"-/healthy", nil)
	if err != nil {
		return err
	}
	for k := range s.Headers {
		req.Header.Add(k, s.Headers[k])
	}
	req = req.WithContext(ctx)
	res, err := s.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error connecting to Victoria Metrics: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("wrong status: %s", res.Status)
	}
	return nil
}

func (s *VMSenor) Update(ctx context.Context) error {
	return fmt.Errorf("not implemented")
}

func (s *VMSenor) Close(ctx context.Context) error {
	// becasuse stateless http connections are being used
	return nil
}
