package vmpusher

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// AddLevelValue will add a water level value of a given station into the internal buffer
func (c *Controller) AddLevelValue(site, station string, lat, long float64, t time.Time, level float64) {
	addValue(c.levels, metricLevelName, site, station, lat, long, t, level)
}

// AddFlowValue will add a water flow value of a given station into the internal buffer
func (c *Controller) AddFlowValue(site, station string, lat, long float64, t time.Time, flow float64) {
	addValue(c.flows, metricFlowName, site, station, lat, long, t, flow)
}

// SendValues will push all the values within the internal buffer to victoria metrics and flush the buffer if successfull
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
