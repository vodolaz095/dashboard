package victoriametrics

import (
	"strconv"
	"time"
)

// {"status":"success","data":{"resultType":"vector","result":[]},"stats":{"seriesFetched": "0","executionTimeMsec":0}}

type rawMetric struct {
	Metric map[string]string `json:"metric"`
	Group  int               `json:"group"`
	Value  []any             `json:"value"`
}

func (rw *rawMetric) hasAllTags(needle map[string]string) (ok bool) {
	var val string
	var present bool
	if len(needle) == 0 {
		return true
	}
	for k := range needle {
		val, present = rw.Metric[k]
		if present && val == needle[k] {
			ok = true
		}
	}
	return ok
}

func (rw *rawMetric) GetLastValue() (val float64, present bool) {
	var ok bool
	var raw string
	var err error
	if len(rw.Value) > 1 {
		raw, ok = rw.Value[1].(string)
		if !ok {
			return 0, false
		}
		val, err = strconv.ParseFloat(raw, 64)
		if err != nil {
			return 0, false
		}
		return val, true
	}
	return 0, false
}

func (rw *rawMetric) GetLastTimestamp() (val time.Time, present bool) {
	var ok bool
	var n float64
	if len(rw.Value) > 0 {
		n, ok = rw.Value[0].(float64)
		if !ok {
			return time.Time{}, false
		}
		return time.UnixMilli(int64(n * 1000)), true
	}
	return time.Time{}, false
}

type rawData struct {
	Status     string             `json:"string"`
	Result     []rawMetric        `json:"result"`
	Stats      map[string]float64 `json:"stats"`
	ResultType string             `json:"resultType"`
}

type rawResponse struct {
	Status string  `json:"status"`
	Data   rawData `json:"data"`
}
