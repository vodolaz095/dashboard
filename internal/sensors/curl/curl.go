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

	"github.com/vodolaz095/dashboard/internal/sensors"
)

type Sensor struct {
	sensors.UnimplementedSensor
	// Client is used to make HTTP connections
	Client *http.Client
	// HttpMethod defines request type being send to remote http(s) endpoint via HTTP protocol
	HttpMethod string `yaml:"http_method" validate:"oneof=GET HEAD POST PUT PATCH DELETE CONNECT OPTIONS TRACE"`
	// Endpoint defines URL where sensor sends request to recieve data
	Endpoint string `yaml:"endpoint" validate:"http_url"`
	// Headers are HTTP request headers being send with any HTTP request
	Headers map[string]string `yaml:"headers"`
	// Body is send with any HTTP request as payload
	Body string `yaml:"body"`
	// JsonPath is used to extract elements from json response of remote endpoint or shell command output using https://jsonpath.com/ syntax
	JsonPath string `yaml:"json_path"`
	// ExpectedStatusCode
	ExpectedStatusCode int `yaml:"expected_status_code"`
}

func (s *Sensor) Init(ctx context.Context) error {
	if s.A == 0 {
		s.A = 1
	}
	if s.HttpMethod == "" {
		s.HttpMethod = http.MethodGet
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
	s.Mutex = &sync.RWMutex{}
	return nil
}

func (s *Sensor) Ping(ctx context.Context) error {
	return nil
}

func (s *Sensor) Close(ctx context.Context) error {
	return nil
}

func (s *Sensor) Update(ctx context.Context) (err error) {
	s.Mutex.Lock()
	defer func() {
		s.Mutex.Unlock()
		if err != nil {
			s.Error = err
		}
	}()
	s.Value = 0
	s.UpdatedAt = time.Now()
	var val float64
	body := bytes.NewBufferString(s.Body)
	req, err := http.NewRequest(s.HttpMethod, s.Endpoint, body)
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
		s.Error = nil
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
	s.Error = nil
	s.Value = res.(float64)
	return nil
}
