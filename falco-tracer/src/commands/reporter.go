package commands

import (
	"fmt"
	"log"
	"os"
	"stats_getter"
)

const (
	jsonFormat   = "json"
	dotFormat    = "dot"
	foldedFormat = "folded"

	jsonFile   = "falco_metrics.json"
	dotFile    = "falco_stacktrace.dot"
	foldedFile = "falco_stacktrace.folded"
)

type reporterData struct {
	falcoBins []string
	falcoargs []string

	outputDirectory string
	outputFormats   []string
	mode            string

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

	fmt.Println("File writer: ", l, "bytes written successfully on ", outPath)
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
