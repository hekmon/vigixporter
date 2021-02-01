package watcher

import (
	"log"
	"time"

	"github.com/hekmon/vigixporter/hubeau"
)

const (
	watchInterval = 10 * time.Minute
)

func (c *Controller) batch() {

	metrics, err := c.source.GetAllObservations(c.ctx, hubeau.ObservationsRequest{
		EntityCode: c.stations,
		Type:       hubeau.ObservationTypeLevelAndFlow,
		StartDate:  c.lastSeen,
		Sort:       hubeau.SortAscending,
		// Timestep:   10,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, metric := range metrics {
		if metric.ObsDate.Equal(c.lastSeen) {
			continue
		}
		switch metric.Type {
		case hubeau.ObservationTypeLevel:
			c.target.AddLevelValue(metric.SiteCode, metric.StationCode, metric.Latitude, metric.Longitude, metric.ObsDate, metric.ObsResultat)
		case hubeau.ObservationTypeFlow:
			c.target.AddFlowValue(metric.SiteCode, metric.StationCode, metric.Latitude, metric.Longitude, metric.ObsDate, metric.ObsResultat)
		default:
			if err != nil {
				log.Fatalf("Unknown type: %s\n", metric.Type)
			}
		}
	}

	err = c.target.SendValues()
	if err != nil {
		log.Fatal(err)
	}
}
