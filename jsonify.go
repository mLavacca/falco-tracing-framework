package main

import (
	"encoding/json"
	"log"
)

func jsonifyFalcoStats(sa StatsAggregator) string {
	jsonStats, err := json.Marshal(sa.Fs)

	if err != nil {
		log.Fatalf("Json marshaling failed, %s", err)
	}

	return string(jsonStats)
}
