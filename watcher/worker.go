package watcher

import (
	"time"

	"github.com/hekmon/vigixporter/hubeau/hydrometrie"
)

const (
	watchInterval = 10 * time.Minute
)

func (c *Controller) worker() {
	// First batch now
	c.batch()
	// Next batchs
	ticker := time.NewTicker(watchInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.batch()
		case <-c.ctx.Done():
			c.logger.Debug("[Watcher] stopping worker")
			return
		}
	}
}

func (c *Controller) batch() {
	// Get metrics
	c.logger.Infof("[Watcher] current batch: getting new metrics...")
	oldestLastSeen := c.getOldestSeen()
	if !oldestLastSeen.IsZero() {
		c.logger.Debugf("[Watcher] current batch: requesting data from the oldest last seen we got: %v", oldestLastSeen)
	}
	metrics, err := c.source.GetAllObservations(c.ctx, hydrometrie.ObservationsRequest{
		EntityCode: c.stations,
		Type:       hydrometrie.ObservationTypeLevelAndFlow,
		StartDate:  oldestLastSeen,
		Sort:       hydrometrie.SortAscending,
	})
	if err != nil {
		c.logger.Errorf("[Watcher] current batch: can't get metrics from hubeau: %s", err)
		return
	}
	if len(metrics) == 0 {
		c.logger.Warning("[Watcher] current batch: we did not retreive any metrics")
		return
	}
	c.logger.Infof("[Watcher] current batch: recovered %d values for %d stations", len(metrics), len(c.stations))
	// Ingerate metrics
	for index, metric := range metrics {
		switch metric.Type {
		case hydrometrie.ObservationTypeLevel:
			if c.isLevelValueKnown(metric.StationCode, metric.ObsDate) {
				c.logger.Debugf("[Watcher] current batch: index %d: level metric has a known date: skipping", index)
				continue
			}
			c.logger.Debugf("[Watcher] current batch: index %d: adding a level metric (station: %s, time: %s, value: %f)",
				index, metric.StationCode, metric.ObsDate, metric.ObsResultat)
			c.target.AddLevelValue(metric.SiteCode, metric.StationCode, metric.Latitude,
				metric.Longitude, metric.ObsDate, metric.ObsResultat)
			c.lastSeenLevelCandidate(metric.StationCode, metric.ObsDate)
		case hydrometrie.ObservationTypeFlow:
			if c.isFlowValueKnown(metric.StationCode, metric.ObsDate) {
				c.logger.Debugf("[Watcher] current batch: index %d: flow metric has a known date: skipping", index)
				continue
			}
			c.logger.Debugf("[Watcher] current batch: index %d: adding a flow metric (station: %s, time: %s, value: %f)",
				index, metric.StationCode, metric.ObsDate, metric.ObsResultat)
			c.target.AddFlowValue(metric.SiteCode, metric.StationCode, metric.Latitude,
				metric.Longitude, metric.ObsDate, metric.ObsResultat)
			c.lastSeenFlowCandidate(metric.StationCode, metric.ObsDate)
		default:
			c.logger.Warningf("[Watcher] current batch: index %d: unknown metric type '%s' has been skipped: %+v",
				index, metric.Type, metric)
		}
	}
	c.logger.Debugf("[Watcher] current batch: updated lastseen:\n\tlevels: %+v\n\tflows: %+v", c.lastSeenLevels, c.lastSeenFlows)
	// Send them to victoria metrics
	nbMetrics, nbValues, err := c.target.Send()
	if err != nil {
		c.logger.Errorf("[Watcher] current batch: can't send metrics to victoria metrics: %s", err)
		return
	}
	if nbMetrics == 0 {
		c.logger.Info("[Watcher] current batch: no metric has been sent")
	} else {
		c.logger.Infof("[Watcher] current batch: successfully sent %d metrics containing %d values", nbMetrics, nbValues)
	}
}

func (c *Controller) getOldestSeen() (oldest time.Time) {
	for _, lastLevelSeen := range c.lastSeenLevels {
		if oldest.IsZero() {
			oldest = lastLevelSeen
		} else if lastLevelSeen.Before(oldest) {
			oldest = lastLevelSeen
		}
	}
	for _, lastFlowSeen := range c.lastSeenFlows {
		if oldest.IsZero() {
			oldest = lastFlowSeen
		} else if lastFlowSeen.Before(oldest) {
			oldest = lastFlowSeen
		}
	}
	return
}

func (c *Controller) isLevelValueKnown(station string, metricTime time.Time) bool {
	return isKnown(c.lastSeenLevels, station, metricTime)
}

func (c *Controller) isFlowValueKnown(station string, metricTime time.Time) bool {
	return isKnown(c.lastSeenFlows, station, metricTime)
}

func isKnown(db map[string]time.Time, id string, value time.Time) (known bool) {
	if lastKnown, found := db[id]; found && !value.After(lastKnown) {
		known = true
	}
	return
}

func (c *Controller) lastSeenLevelCandidate(station string, candidate time.Time) {
	lastSeenCandidate(c.lastSeenLevels, station, candidate)
}

func (c *Controller) lastSeenFlowCandidate(station string, candidate time.Time) {
	lastSeenCandidate(c.lastSeenFlows, station, candidate)
}

func lastSeenCandidate(db map[string]time.Time, id string, candidate time.Time) {
	if ref, found := db[id]; !found {
		db[id] = candidate
	} else if candidate.After(ref) {
		db[id] = candidate
	}
}
