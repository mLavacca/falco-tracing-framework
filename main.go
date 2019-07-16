package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func main() {

	var wg sync.WaitGroup

	t, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal("error, time parameter missing")
	}

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	falcoTracer := NewFalcoTracer()

	falcoTracer.setupConnection()

	falcoTracer.loadRulesFromFalco()

	falcoTracer.flushFalcoData()

	wg.Add(1)
	go falcoTracer.loadStatsFromFalco(time.Duration(t), &wg)

	<-sigs

	falcoTracer.exitFlag = true

	wg.Wait()

	jsonStats, err := falcoTracer.MarshalJSON()
	if err != nil {
		log.Fatal("error in object marshaling")
	}

	fmt.Println(string(jsonStats))
}
