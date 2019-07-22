package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"os/user"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func main() {

	var wg sync.WaitGroup

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	DesktopDir := user.HomeDir + "/Desktop/"

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

	f, err := os.Create(DesktopDir + "tracer_data.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	l, err := f.Write(jsonStats)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
