package victoriametrics

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

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
	// See https://docs.victoriametrics.com/url-examples/#apiv1query
	var raw rawResponse
	u, err := url.Parse(s.Endpoint)
	if err != nil {
		return fmt.Errorf("malformed endpoint: %w", err)
	}
	defer func() {
		if err != nil {
			s.Mutex.Lock()
			s.Error = err
			s.Value = 0
			s.UpdatedAt = time.Now()
			s.Mutex.Unlock()
		}
	}()
	u.Path += "prometheus/api/v1/query"
	params := url.Values{}
	params.Set("query", s.Query)
	params.Set("time", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("step", time.Second.String())
	deadline, present := ctx.Deadline()
	if present {
		params.Set("timeout", time.Until(deadline).String())
	}
	u.RawQuery = params.Encode()
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req = req.WithContext(ctx)
	for k := range s.Headers {
		req.Header.Set(k, s.Headers[k])
	}
	resp, err := s.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request to %s: %w", req.URL.String(), err)
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("wrong response: %s", resp.Status)
	}
	err = json.NewDecoder(resp.Body).Decode(&raw)
	if err != nil {
		return fmt.Errorf("error decoding body: %w", err)
	}
	if raw.Status != "success" {
		return fmt.Errorf("wrong query status: %s", raw.Status)
	}
	for i := range raw.Data.Result {
		if raw.Data.Result[i].hasAllTags(s.Filter) {
			val, found1 := raw.Data.Result[i].GetLastValue()
			when, found2 := raw.Data.Result[i].GetLastTimestamp()
			if found1 && found2 {
				s.Mutex.Lock()
				log.Debug().Msgf("VM: updating sensor %s to %.4f on %s", s.Name, val, when.Format("15:04:05"))
				s.Value = val
				s.UpdatedAt = when
				s.Mutex.Unlock()
				return nil
			}
			s.Mutex.Lock()
			s.Value = 0
			s.UpdatedAt = time.Now()
			s.Error = fmt.Errorf("no data")
			s.Mutex.Unlock()
		}
	}
	return nil
}

func (s *VMSenor) Close(ctx context.Context) error {
	// because stateless http connections are being used
	return nil
}
