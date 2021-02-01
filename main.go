package main

import (
	"context"

	"github.com/hekmon/vigixporter/hubeau"
	"github.com/hekmon/vigixporter/watcher"
)

func main() {
	listOfStations := []string{hubeau.StationParis, hubeau.StationAlfortville, hubeau.StationCreteil}

	watcher, err := watcher.New(context.TODO(), watcher.Config{
		Stations: listOfStations,
	})
}
