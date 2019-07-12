package main

import (
	"fmt"
	"time"
)

func main() {

	// placeholder to modify in the future (get parameter from command line or something else)
	t := time.Duration(4)

	ch := make(chan StatsAggregator)

	falcoTracer := NewFalcoTracer()

	falcoTracer.setupConnection()

	falcoTracer.loadRulesFromFalco()

	go falcoTracer.loadStatsFromFalco(t, ch)

	for {
		sa := <-ch
		jsonStats := jsonifyFalcoStats(sa)

		fmt.Println(jsonStats)
	}
}
