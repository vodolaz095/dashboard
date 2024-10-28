package model

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

// Sensor is data transfer object for http, json and metrics visualization
type Sensor struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Description string            `json:"description"`
	Link        string            `json:"link"`
	Minimum     float64           `json:"minimum"`
	Maximum     float64           `json:"maximum"`
	Value       float64           `json:"value"`
	Error       string            `json:"error"`
	Tags        map[string]string `json:"tags"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}

// String formats sensor readings in a way prometheus understands it
func (s *Sensor) String() string {
	// https://prometheus.io/docs/instrumenting/exposition_formats/#text-format-example
	/*
		# HELP http_requests_total The total number of HTTP requests.
		# TYPE http_requests_total gauge
		http_requests_total{method="post",code="200"} 1027 1395066363000
	*/
	labels := bytes.NewBufferString("")
	tags := make([]string, len(s.Tags))
	i := 0
	for k := range s.Tags {
		tags[i] = fmt.Sprintf("%q=%q", k, s.Tags[k])
		i++
	}
	buh := bytes.NewBufferString("")
	if len(s.Tags) > 0 {
		labels.WriteString("{")
		labels.WriteString(strings.Join(tags, ","))
		labels.WriteString("}")
	}
	fmt.Fprintln(buh, "# HELP", s.Name, s.Description, s.Link)
	fmt.Fprintln(buh, "# TYPE", s.Name, "gauge")
	fmt.Fprintf(buh, "%s%s %v %v\n", s.Name, labels.String(), s.Value, s.UpdatedAt.Unix())
	return buh.String()
}

// GetStatus returns sensor status - ok, low, high
func (s *Sensor) GetStatus() string {
	if s.Minimum == 0 && s.Maximum == 0 {
		return StatusOK
	}
	if s.Minimum > s.Value {
		return StatusLow
	}
	if s.Maximum < s.Value {
		return StatusHigh
	}
	return StatusOK
}
