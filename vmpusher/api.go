package vmpusher

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func (c *Controller) AddLevelValue(site, station string, lat, long float64, t time.Time, level float64) {
	addValue(c.levels, metricLevelName, site, station, lat, long, t, level)
}

func (c *Controller) AddFlowValue(site, station string, lat, long float64, t time.Time, flow float64) {
	addValue(c.flows, metricFlowName, site, station, lat, long, t, flow)
}

func (c *Controller) SendValues() (err error) {
	buffer := new(strings.Builder)
	encoder := json.NewEncoder(buffer)
	// write levels
	for station, levelmetric := range c.levels {
		if err = encoder.Encode(levelmetric); err != nil {
			return fmt.Errorf("can't encode level metrics for station '%s': %w", station, err)
		}
	}
	// write flows
	for station, flowmetric := range c.flows {
		if err = encoder.Encode(flowmetric); err != nil {
			return fmt.Errorf("can't encode level metrics for station '%s': %w", station, err)
		}
	}
	// send payload
	fmt.Println(buffer.String())
	// payload successfully sent
	clearValues(c.levels)
	clearValues(c.flows)
	return
}
