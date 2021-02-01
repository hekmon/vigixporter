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

type jsonlMetric struct {
	Metric     map[string]string `json:"metric"`
	Values     []float64         `json:"values"`
	Timestamps []int64           `json:"timestamps"`
}

func newjsonlMetric(name string, labels map[string]string) jsonlMetric {
	labels["__name__"] = name
	return jsonlMetric{Metric: labels}
}

func addValue(metrics map[string]jsonlMetric, metricType, site, station string, lat, long float64, time time.Time, value float64) {
	var (
		metric jsonlMetric
		found  bool
	)
	if metric, found = metrics[station]; !found {
		metric = newjsonlMetric(fmt.Sprintf("%s_%s_%s", metricNamespace, metricSubsystem, metricType), map[string]string{
			labelSite:    site,
			labelStation: station,
			labelLat:     fmt.Sprintf("%f", lat),
			labelLong:    fmt.Sprintf("%f", long),
		})
		metrics[station] = metric
	}
	metric.Timestamps = append(metric.Timestamps, time.Unix())
	metric.Values = append(metric.Values, value)
}

func clearValues(metrics map[string]jsonlMetric) {
	for _, metric := range metrics {
		metric.Timestamps = make([]int64, 0, 1)
		metric.Values = make([]float64, 0, 1)
	}
}
