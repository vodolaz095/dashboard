package shell

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/oliveagle/jsonpath"
	"github.com/vodolaz095/dashboard/sensors"
)

type Sensor struct {
	sensors.UnimplementedSensor
	val       float64
	updatedAt time.Time
}

func (s *Sensor) Init(ctx context.Context) error {
	args := strings.Split(s.DatabaseConnectionString, " ")
	stat, err := os.Stat(args[0])
	if err != nil {
		return err
	}
	if stat.Mode()&0111 == 0 {
		return fmt.Errorf("file %s is not executable", s.DatabaseConnectionString)
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

func (s *Sensor) Update(ctx context.Context, _ float64) (err error) {
	var val float64
	args := strings.Split(s.DatabaseConnectionString, " ")
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	raw, err := cmd.Output()
	if err != nil {
		return
	}
	// no processing script output
	if s.Query == "" {
		val, err = strconv.ParseFloat(strings.TrimSpace(string(raw)), 64)
		if err != nil {
			return
		}
		s.val = val
		s.updatedAt = time.Now()
		return nil
	}
	// command returned json we need to execute jsonpath query against
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
