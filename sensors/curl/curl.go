package curl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/oliveagle/jsonpath"
	"github.com/vodolaz095/dashboard/sensors"
)

type Sensor struct {
	sensors.UnimplementedSensor
	mu                 *sync.Mutex
	Headers            map[string]string
	Method             string
	Body               string
	ExpectedStatusCode int
	Client             *http.Client
}

func (s *Sensor) Init(ctx context.Context) error {
	if s.A == 0 {
		s.A = 1
	}
	if s.Method == "" {
		s.Method = http.MethodGet
	}
	if s.ExpectedStatusCode == 0 {
		s.ExpectedStatusCode = http.StatusOK
	}
	if s.Headers == nil {
		s.Headers = make(map[string]string, 0)
	}
	if s.Client == nil {
		s.Client = http.DefaultClient
	}
	s.mu = &sync.Mutex{}
	return nil
}

func (s *Sensor) Ping(ctx context.Context) error {
	return nil
}

func (s *Sensor) Close(ctx context.Context) error {
	return nil
}

func (s *Sensor) Update(ctx context.Context) (err error) {
	s.mu.Lock()
	defer func() {
		s.mu.Unlock()
		if err != nil {
			s.Error = err
		}
	}()
	s.Value = 0
	s.Error = nil
	s.UpdatedAt = time.Now()
	var val float64
	body := bytes.NewBufferString(s.Body)
	req, err := http.NewRequest(s.Method, s.Endpoint, body)
	if err != nil {
		return
	}
	for k, v := range s.Headers {
		req.Header.Add(k, v)
	}
	req = req.WithContext(ctx)
	resp, err := s.Client.Do(req)
	if err != nil {
		return
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if resp.StatusCode != s.ExpectedStatusCode {
		s.Error = fmt.Errorf("unexpected status %v %s", resp.StatusCode, resp.Status)
		return s.Error
	}
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		s.Error = err
		return
	}
	if s.JsonPath == "" {
		val, err = strconv.ParseFloat(strings.TrimSpace(string(raw)), 64)
		if err != nil {
			s.Error = err
			return
		}
		s.Value = val
		s.UpdatedAt = time.Now()
		return
	}

	var data interface{}
	err = json.Unmarshal(raw, &data)
	if err != nil {
		s.Error = err
		return
	}
	res, err := jsonpath.JsonPathLookup(data, s.JsonPath)
	if err != nil {
		s.Error = err
		return
	}
	s.Value = res.(float64)
	return nil
}
