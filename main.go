package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hekmon/vigixporter/hubeau"
)

const (
	stationParis       = "F700000103"
	stationAlfortville = "F490000104"
	stationCreteil     = "F664000404"
)

func main() {
	glouglou := hubeau.New()

	answer, err := glouglou.GetObservations(context.Background(), hubeau.ObservationsRequest{
		EntityCode: []string{stationParis, stationAlfortville, stationCreteil},
		Type:       hubeau.ObservationTypeHeightAndSpeed,
		// Size:       hubeau.RequestMaxSize,
		Size: 20,
		Sort: hubeau.SortDescending,
		// Timestep:   10,
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", answer)
}
