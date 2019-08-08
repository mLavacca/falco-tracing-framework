package main

import (
	"log"
	"os"
)

var functionsSlice = map[int]interface{}{
	0: writeBelowRoot,
	1: writeBelowEtc,
}

var rollbacksSlice = map[int]interface{}{
	1: writeBelowEtcRollback,
}

/*
 * id = 0
 */
func writeBelowRoot() {
	path := "/sbin/iptables"

	f, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

/*
 * id = 1
 */
func writeBelowEtc() {
	path := "/etc/falco_tester_file"

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func writeBelowEtcRollback() {
	path := "/etc/falco_tester_file"

	err := os.Remove(path)
	if err != nil {
		log.Fatal(err)
	}
}
