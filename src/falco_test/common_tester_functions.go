package falco_test

import (
	"log"
	"os"
)

func writeFile(path string) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func deleteFile(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Fatal(err)
	}
}
