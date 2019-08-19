package commands

import (
	"fmt"
	"log"
	"os"
	"stats_getter"
)

type reporterData struct {
	falcoBin  string
	falcoargs []string

	outputFile       string
	outputFoldedFile string
	mode             string

	falcoTracer *stats_getter.FalcoTracer
}

func writeMetricsOnFile(data []byte, outPath string) {

	f, err := os.Create(outPath)
	if err != nil {
		log.Fatalln("Error in file creation", outPath, err)
	}

	l, err := f.Write(data)
	if err != nil {
		log.Fatalln("Error in json write", err)
	}

	fmt.Println("json output: ", l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
