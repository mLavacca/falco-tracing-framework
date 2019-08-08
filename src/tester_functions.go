package main

import (
	"log"
	"os"
)

const functionsSlice []interface{}{
	writeBelowRoot
	writeBelowEtc
}

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
