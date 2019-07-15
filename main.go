package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {

	t, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal("error, time parameter missing")
	}

	ch := make(chan StatsAggregator)

	falcoTracer := NewFalcoTracer()

	falcoTracer.setupConnection()

	falcoTracer.loadRulesFromFalco()

	falcoTracer.flushFalcoData()

	go falcoTracer.loadStatsFromFalco(time.Duration(t), ch)

	for {
		sa := <-ch
		jsonStats, err := sa.MarshalJSON()
		if err != nil {
			log.Fatal("error in object marshaling")
		}

		fmt.Print(string(jsonStats), "\n\n")
	}
}
