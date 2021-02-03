package vmpusher

import (
	"fmt"
	"time"
)

const (
	// Name composition
	metricNamespace = "vigixporter"
	metricSubsystem = "water"
	metricLevelName = "level"
	metricFlowName  = "flow"
	// labels
	labelSite    = "site_code"
	labelStation = "station_code"
	labelLat     = "latitude"
	labelLong    = "longitude"
)

// JSONLineMetric represents a given metric with N values in the JSON line (streaming) format
type JSONLineMetric struct {
	Metric     map[string]string `json:"metric"`
	Values     []float64         `json:"values"`
	Timestamps []int64           `json:"timestamps"`
}

func newJSONLineMetric(name string, labels map[string]string) JSONLineMetric {
	labels["__name__"] = name
	return JSONLineMetric{Metric: labels}
}

func addValue(metrics map[string]JSONLineMetric, metricType, site, station string, lat, long float64, time time.Time, value float64) {
	var (
		metric JSONLineMetric
		found  bool
	)
	if metric, found = metrics[station]; !found {
		metric = newJSONLineMetric(fmt.Sprintf("%s_%s_%s", metricNamespace, metricSubsystem, metricType), map[string]string{
			labelSite:    site,
			labelStation: station,
			labelLat:     fmt.Sprintf("%f", lat),
			labelLong:    fmt.Sprintf("%f", long),
		})
	}
	metric.Timestamps = append(metric.Timestamps, time.Unix()*1000)
	metric.Values = append(metric.Values, value)
	metrics[station] = metric
}

func clearValues(metrics map[string]JSONLineMetric) {
	for station, metric := range metrics {
		metric.Timestamps = make([]int64, 0, 1)
		metric.Values = make([]float64, 0, 1)
		metrics[station] = metric
	}
}
