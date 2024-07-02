package shell

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/oliveagle/jsonpath"
	"github.com/vodolaz095/dashboard/sensors"
)

type Sensor struct {
	sensors.UnimplementedSensor

	// Command is shell command being executed by shell sensor
	Command string `yaml:"command"`
	// Environment is POSIX environment used by shell sensor to execute commands into
	Environment map[string]string `yaml:"headers"`
	// JsonPath is used to extract elements from json response of remote endpoint or shell command output using https://jsonpath.com/ syntax
	JsonPath string `yaml:"json_path"`
}

func (s *Sensor) Init(ctx context.Context) (err error) {
	args := strings.Split(s.Command, " ")
	s.Mutex = &sync.RWMutex{}
	if s.A == 0 {
		s.A = 1
	}
	_, err = exec.LookPath(args[0])
	return err
}

func (s *Sensor) Ping(ctx context.Context) error {
	return nil
}

func (s *Sensor) Close(ctx context.Context) error {
	return nil
}

func (s *Sensor) Update(ctx context.Context) (err error) {
	var val float64
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.UpdatedAt = time.Now()
	s.Error = nil
	s.Value = 0
	args := strings.Split(s.Command, " ")
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	for k := range s.Environment {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%q", k, s.Environment[k]))
	}
	raw, err := cmd.Output()
	if err != nil {
		s.Error = err
		return
	}
	// no processing script output
	if s.JsonPath == "" {
		val, err = strconv.ParseFloat(strings.TrimSpace(string(raw)), 64)
		if err != nil {
			s.Error = err
			return
		}
		s.Value = val
		s.UpdatedAt = time.Now()
		return nil
	}
	// command returned json we need to execute jsonpath query against
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
