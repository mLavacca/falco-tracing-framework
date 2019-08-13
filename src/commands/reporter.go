package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"stats_getter"
)

type reporterData struct {
	falcoBin  string
	falcoargs []string

	outputFile string
	mode       string

	falcoTracer *stats_getter.FalcoTracer
}

func writeMetricsOnFile(jsonStats []byte, outDir string) {
	outpath := filepath.Join(outDir, "tracer_data.json")

	f, err := os.Create(outpath)
	if err != nil {
		fmt.Println(err)
		return
	}

	l, err := f.Write(jsonStats)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("json output: ", l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
