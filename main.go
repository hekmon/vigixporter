package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hekmon/vigixporter/hubeau"
)

func main() {
	glouglou := hubeau.New()

	answer, err := glouglou.GetObservations(context.Background(), hubeau.ObservationsRequest{
		EntityCode: []string{hubeau.StationAlfortville},
		Type:       hubeau.ObservationTypeHeight,
		StartDate:  time.Now().Add(24 * time.Hour * -1),
		// EndDate:    time.Now(),
		// Size:       hubeau.RequestMaxSize,
		Size: hubeau.RequestMaxSize,
		Sort: hubeau.SortAscending,
		// Timestep:   10,
	})
	if err != nil {
		log.Fatal(err)
	}
	printResults(answer.Data)

	lastdate := answer.Data[len(answer.Data)-1].ObsDate

	answer, err = glouglou.GetObservations(context.Background(), hubeau.ObservationsRequest{
		EntityCode: []string{hubeau.StationAlfortville},
		Type:       hubeau.ObservationTypeHeight,
		StartDate:  lastdate,
		// EndDate:    time.Now(),
		// Size:       hubeau.RequestMaxSize,
		Size: hubeau.RequestMaxSize,
		Sort: hubeau.SortAscending,
		// Timestep:   10,
	})
	if err != nil {
		log.Fatal(err)
	}
	printResults(answer.Data)
}

func printResults(data []hubeau.Observation) {
	// fmt.Printf("%#v\n", answer)
	for _, obs := range data {
		fmt.Printf("%s\t%s", obs.ObsDate, obs.StationCode)
		switch obs.Type {
		case hubeau.ObservationTypeHeight:
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
