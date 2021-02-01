package main

import (
	"context"
	"os"

	"github.com/hekmon/hllogger"
	"github.com/hekmon/vigixporter/hubeau"
	"github.com/hekmon/vigixporter/watcher"
)

var (
	logger *hllogger.HlLogger
)

func main() {
	listOfStations := []string{hubeau.StationParis, hubeau.StationAlfortville, hubeau.StationCreteil}

	logger = hllogger.New(os.Stderr, &hllogger.Config{
		LogLevel: hllogger.Debug,
	})

	watcher, err := watcher.New(context.TODO(), watcher.Config{
		Stations: listOfStations,
		Logger:   logger,
	})
}
