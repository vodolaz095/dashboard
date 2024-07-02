package file

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/oliveagle/jsonpath"
	"github.com/vodolaz095/dashboard/sensors"
)

type Sensor struct {
	sensors.UnimplementedSensor
	// PathToReadingsFile sets path to file we are reading periodically
	PathToReadingsFile string
	// JsonPath is used to extract elements from json response of remote endpoint or shell command output using https://jsonpath.com/ syntax
	JsonPath string `yaml:"json_path"`
}

func (s *Sensor) Init(_ context.Context) error {
	s.Mutex = &sync.RWMutex{}
	if s.A == 0 {
		s.A = 1
	}
	return nil
}

func (s *Sensor) Ping(_ context.Context) error {
	return nil
}

func (s *Sensor) Close(_ context.Context) error {
	return nil
}

func (s *Sensor) Update(ctx context.Context) (err error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.UpdatedAt = time.Now()
	s.Value = 0
	s.Error = nil
	var val float64
	raw, err := os.ReadFile(s.PathToReadingsFile)
	if err != nil {
		return err
	}

	// readings file contains plain old raw value
	if s.JsonPath == "" {
		val, err = strconv.ParseFloat(strings.TrimSpace(string(raw)), 64)
		if err != nil {
			s.Error = err
			return
		}
		s.Value = val
		return nil
	}
	// readings file contains json we need to execute jsonpath query against
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
