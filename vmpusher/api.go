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

func (c *Controller) ClearValues() {
	clearValues(c.levels)
	clearValues(c.flows)
}

func (c *Controller) SendValues() (err error) {
	buffer := new(strings.Builder)
	encoder := json.NewEncoder(buffer)
	// write levels
	for station, levelmetric := range c.levels {
		if err = encoder.Encode(levelmetric); err != nil {
			return fmt.Errorf("can't encode level metrics for station '%s': %w", station, err)
		}
		fmt.Fprintln(buffer)
	}
	// write flows
	for station, flowmetric := range c.flows {
		if err = encoder.Encode(flowmetric); err != nil {
			return fmt.Errorf("can't encode level metrics for station '%s': %w", station, err)
		}
		fmt.Fprintln(buffer)
	}
	// payload ready
	fmt.Println(buffer.String())
	return
}
