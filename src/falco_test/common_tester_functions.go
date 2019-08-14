package falco_test

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
)

func openFile(path string, flags int) {
	f, err := os.OpenFile(path, flags, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	if err := f.Close(); err != nil {
		log.Fatalln(err)
	}
}

func deleteFile(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Fatalln(err)
	}
}

func createDir(path string) {
	err := os.Mkdir(path, 0666)
	if err != nil {
		log.Fatalln(err)
	}
}

func removeDir(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Fatalln(err)
	}
}

func getUsername() (string, error) {
	files, err := ioutil.ReadDir("/home")
	if err != nil {
		log.Fatalln(err)
	}

	for _, f := range files {
		if f.Name() != "root" {
			return f.Name(), nil
		}
	}
	return "", errors.New("error in finding users")
}
