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
	"time"

	"github.com/oliveagle/jsonpath"
	"github.com/vodolaz095/dashboard/sensors"
)

type Sensor struct {
	sensors.UnimplementedSensor
	Headers            map[string]string
	Method             string
	Body               string
	ExpectedStatusCode int
	Client             *http.Client
	val                float64
	updatedAt          time.Time
}

func (s *Sensor) Init(ctx context.Context) error {
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
	return nil
}

func (s *Sensor) Ping(ctx context.Context) error {
	return nil
}

func (s *Sensor) Close(ctx context.Context) error {
	return nil
}

func (s *Sensor) Value() float64 {
	return s.val
}

func (s *Sensor) Update(ctx context.Context) (err error) {
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
		return fmt.Errorf("unexpected status %v %s", resp.StatusCode, resp.Status)
	}
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if s.Query == "" {
		val, err = strconv.ParseFloat(strings.TrimSpace(string(raw)), 64)
		if err != nil {
			return
		}
		s.val = val
		s.updatedAt = time.Now()
		return
	}

	var data interface{}
	err = json.Unmarshal(raw, &data)
	if err != nil {
		return
	}
	res, err := jsonpath.JsonPathLookup(data, s.Query)
	if err != nil {
		return
	}
	s.val = res.(float64)
	s.updatedAt = time.Now()
	return nil
}

func (s *Sensor) UpdatedAt() time.Time {
	return s.updatedAt
}
