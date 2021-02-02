package watcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/hekmon/vigixporter/vmpusher"
)

const (
	stateFile = "vigixporter_state.json"
)

type state struct {
	LevelsBuffer   map[string]vmpusher.JSONLineMetric `json:"levels_buffer"`
	FlowsBuffer    map[string]vmpusher.JSONLineMetric `json:"flows_buffer"`
	LastSeenLevels map[string]time.Time               `json:"levels_lastseen"`
	LastSeenFlows  map[string]time.Time               `json:"flows_lastseen"`
}

func loadState() (s state, err error) {
	// handle file descriptor
	fd, err := os.Open(stateFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// file does not exist, let's start fresh !
			s.LevelsBuffer = make(map[string]vmpusher.JSONLineMetric)
			s.FlowsBuffer = make(map[string]vmpusher.JSONLineMetric)
			s.LastSeenLevels = make(map[string]time.Time)
			s.LastSeenFlows = make(map[string]time.Time)
			err = nil
			return
		}
		// File may exists but we can not open it
		err = fmt.Errorf("can't open %s state file: %w", stateFile, err)
		return
	}
	defer fd.Close()
	// handle content
	if err = json.NewDecoder(fd).Decode(&s); err != nil {
		err = fmt.Errorf("can't parse %s state file: %w", stateFile, err)
		return
	}
	return
}

func saveState(s state) (err error) {
	fd, err := os.OpenFile(stateFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0640)
	if err != nil {
		err = fmt.Errorf("can't open %s state file: %w", stateFile, err)
		return
	}
	defer fd.Close()
	// handle content
	if err = json.NewEncoder(fd).Encode(s); err != nil {
		err = fmt.Errorf("can't write to state file %s: %w", stateFile, err)
		return
	}
	return
}
