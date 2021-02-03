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

// GetLevelsBuffer returns the current levels buffer
func (c *Controller) GetLevelsBuffer() (levels map[string]JSONLineMetric) {
	return c.levels
}

// GetFlowsBuffer returns the current flows buffer
func (c *Controller) GetFlowsBuffer() (levels map[string]JSONLineMetric) {
	return c.flows
}

// Send will push all the values within the internal buffer to victoria metrics and flush the buffer if successfull
func (c *Controller) Send() (nbMetrics, nbValues int, err error) {
	if len(c.flows) == 0 && len(c.levels) == 0 {
		return
	}
	// marshall buffers into jsonl payload
	payload, err := c.preparePayload()
	if err != nil {
		err = fmt.Errorf("can't marshall internal buffers as JSON line payload: %w", err)
		return
	}
	// send payload
	if err = c.push(payload.String()); err != nil {
		err = fmt.Errorf("failed to push the metrics to victoria metrics server: %w", err)
		return
	}
	// compute stats
	for _, levelMetric := range c.levels {
		if len(levelMetric.Values) != 0 {
			nbMetrics++
			nbValues += len(levelMetric.Values)
		}
	}
	for _, flowMetric := range c.flows {
		if len(flowMetric.Values) != 0 {
			nbMetrics++
			nbValues += len(flowMetric.Values)
		}
	}
	// cleanup
	c.levels = make(map[string]JSONLineMetric, len(c.levels))
	c.flows = make(map[string]JSONLineMetric, len(c.flows))
	return
}

func (c *Controller) preparePayload() (payload strings.Builder, err error) {
	encoder := json.NewEncoder(&payload)
	// write levels
	for station, levelmetric := range c.levels {
		if len(levelmetric.Values) != 0 {
			if err = encoder.Encode(levelmetric); err != nil {
				err = fmt.Errorf("can't encode level metrics for station '%s': %w", station, err)
				return
			}
		}
	}
	// write flows
	for station, flowmetric := range c.flows {
		if len(flowmetric.Values) != 0 {
			if err = encoder.Encode(flowmetric); err != nil {
				err = fmt.Errorf("can't encode flow metrics for station '%s': %w", station, err)
				return
			}
		}
	}
	return
}
