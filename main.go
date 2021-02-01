package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hekmon/vigixporter/hubeau"
	"github.com/hekmon/vigixporter/vmpusher"
)

func main() {
	listOfStations := []string{hubeau.StationParis, hubeau.StationAlfortville, hubeau.StationCreteil}

	glouglou := hubeau.New()

	metrics, err := glouglou.GetAllObservations(context.Background(), hubeau.ObservationsRequest{
		EntityCode: listOfStations,
		Type:       hubeau.ObservationTypeLevelAndFlow,
		// StartDate:  time.Now().Add(24 * time.Hour * -1),
		// EndDate:    time.Now(),
		Sort: hubeau.SortAscending,
		// Timestep:   10,
	})
	if err != nil {
		log.Fatal(err)
	}
	// printResults(answer.Data)

	pusher := vmpusher.New()

	for _, metric := range metrics {
		switch metric.Type {
		case hubeau.ObservationTypeLevel:
			pusher.AddLevelValue(metric.SiteCode, metric.StationCode, metric.Latitude, metric.Longitude, metric.ObsDate, metric.ObsResultat)
		case hubeau.ObservationTypeFlow:
			pusher.AddFlowValue(metric.SiteCode, metric.StationCode, metric.Latitude, metric.Longitude, metric.ObsDate, metric.ObsResultat)
		default:
			if err != nil {
				log.Fatalf("Unknown type: %s\n", metric.Type)
			}
		}
	}

	err = pusher.SendValues()
	if err != nil {
		log.Fatal(err)
	}
}

func printResults(data []hubeau.Observation) {
	// fmt.Printf("%#v\n", answer)
	for _, obs := range data {
		fmt.Printf("%s\t%s", obs.ObsDate, obs.StationCode)
		switch obs.Type {
		case hubeau.ObservationTypeLevel:
			fmt.Printf("\t%.2fm", obs.ObsResultat/1000)
		case hubeau.ObservationTypeFlow:
			fmt.Printf("\t%.2fm3/s", obs.ObsResultat/1000)
		default:
			fmt.Printf("\t%f(unknown)", obs.ObsResultat)
		}
		if obs.SerieStatus != 4 {
			fmt.Printf("\t(unknown status: %d)", obs.SerieStatus)
		}
		fmt.Println()
	}
	fmt.Printf("Printed %d results\n", len(data))
}
