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

func (c *Controller) GetBuffers() (levels, flows map[string]JSONLineMetric) {
	return c.levels, c.flows
}

// SendValues will push all the values within the internal buffer to victoria metrics and flush the buffer if successfull
func (c *Controller) SendValues() (err error) {
	// marshall buffers into jsonl payload
	payload, err := c.preparePayload()
	if err != nil {
		err = fmt.Errorf("can't marshall internal buffers as JSON line payload: %w", err)
		return
	}
	// send payload
	fmt.Println(payload.String())
	// payload successfully sent
	clearValues(c.levels)
	clearValues(c.flows)
	return
}

func (c *Controller) preparePayload() (payload strings.Builder, err error) {
	encoder := json.NewEncoder(&payload)
	// write levels
	for station, levelmetric := range c.levels {
		if err = encoder.Encode(levelmetric); err != nil {
			err = fmt.Errorf("can't encode level metrics for station '%s': %w", station, err)
			return
		}
	}
	// write flows
	for station, flowmetric := range c.flows {
		if err = encoder.Encode(flowmetric); err != nil {
			err = fmt.Errorf("can't encode flow metrics for station '%s': %w", station, err)
			return
		}
	}
	return
}
