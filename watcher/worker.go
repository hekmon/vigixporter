package watcher

import (
	"time"

	"github.com/hekmon/vigixporter/hubeau"
)

const (
	watchInterval = 10 * time.Minute
)

func (c *Controller) batch() {
	// Get metrics
	metrics, err := c.source.GetAllObservations(c.ctx, hubeau.ObservationsRequest{
		EntityCode: c.stations,
		Type:       hubeau.ObservationTypeLevelAndFlow,
		StartDate:  c.lastSeen,
		Sort:       hubeau.SortAscending,
		// Timestep:   10,
	})
	if err != nil {
		c.logger.Errorf("[Watcher] current batch: can't get metrics from hubeau: %s", err)
		return
	}
	if len(metrics) == 0 {
		c.logger.Warning("[Watcher] current batch: we did not retreive any metrics")
		return
	}
	c.logger.Infof("[Watcher] current batch: recovered %d metrics for %d stations", len(metrics), len(c.stations))
	// Ingerate metrics
	var oldest time.Time
	for index, metric := range metrics {
		if metric.ObsDate.Equal(c.lastSeen) {
			c.logger.Debugf("[Watcher] current batch: index %d: metric has a known date: skipping", index)
			continue
		}
		switch metric.Type {
		case hubeau.ObservationTypeLevel:
			c.logger.Debugf("[Watcher] current batch: index %d: adding a level metric (station: %d, time: %s, value: %f)",
				index, metric.StationCode, metric.ObsDate, metric.ObsResultat)
			c.target.AddLevelValue(metric.SiteCode, metric.StationCode, metric.Latitude,
				metric.Longitude, metric.ObsDate, metric.ObsResultat)
		case hubeau.ObservationTypeFlow:
			c.logger.Debugf("[Watcher] current batch: index %d: adding a flow metric (station: %d, time: %s, value: %f)",
				index, metric.StationCode, metric.ObsDate, metric.ObsResultat)
			c.target.AddFlowValue(metric.SiteCode, metric.StationCode, metric.Latitude,
				metric.Longitude, metric.ObsDate, metric.ObsResultat)
		default:
			c.logger.Warningf("[Watcher] current batch: index %d: unknown metric type '%s' has been skipped: %+v",
				index, metric.Type, metric)
		}
		if metric.ObsDate.After(oldest) {
			oldest = metric.ObsDate
			c.logger.Debugf("[Watcher] current batch: oldest date seen so far: %s", oldest)
		}
	}
	c.lastSeen = oldest
	// Send them to victoria metrics
	nbMetrics, err := c.target.Send()
	if err != nil {
		c.logger.Errorf("[Watcher] current batch: can't send metrics to victoria metrics: %s", err)
		return
	}
	if nbMetrics == 0 {
		c.logger.Info("[Watcher] current batch: no metric has been sent")
	} else {
		c.logger.Infof("[Watcher] current batch: successfully sent %d metrics", nbMetrics)
	}
}
