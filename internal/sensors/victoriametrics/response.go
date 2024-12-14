package victoriametrics

import "time"

// {"status":"success","data":{"resultType":"vector","result":[]},"stats":{"seriesFetched": "0","executionTimeMsec":0}}

type rawMetric struct {
	Metric map[string]string `json:"metric"`
	Group  int               `json:"group"`
	Value  []any             `json:"value"`
}

func (rw *rawMetric) hasAllTags(needle map[string]string) (ok bool) {
	var val string
	var present bool

	for k := range needle {
		val, present = rw.Metric[k]
		if present && val != needle[k] {
			return false
		}
		ok = true
	}

	return ok
}

func (rw *rawMetric) GetLastValue() (val float64, present bool) {
	var ok bool
	if len(rw.Value) > 1 {
		val, ok = rw.Value[1].(float64)
		if !ok {
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
