package stats_getter

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"syscall"
)

const rcvStatsSignal = 34
const flushDataSignal = 35
const rcvRulesNamesSignal = 36

type FalcoGateway struct {
	mode          string
	falcoPid      int
	inputFileName string
	pipeFile      *os.File

	sigRcvSummary    syscall.Signal
	sigRcvRulesNames syscall.Signal
	sigFlushData     syscall.Signal

	pipeReader *bufio.Reader
}

func (f *FalcoGateway) configureSignals() {

	f.falcoPid = 0

	f.sigRcvSummary = syscall.Signal(rcvStatsSignal)
	f.sigRcvRulesNames = syscall.Signal(rcvRulesNamesSignal)
	f.sigFlushData = syscall.Signal(flushDataSignal)
}

func NewFalcoGateway(mode string) *FalcoGateway {

	fg := new(FalcoGateway)
	fg.mode = mode

	var inputPath string

	if mode == "online" {
		inputPath = "/tmp/TO_DO"
		fg.configureSignals()
	}

	if mode == "offline" {
		inputPath = "/tmp/falco_tracer_file"
	}

	fg.inputFileName = inputPath
	fg.openPipe()

	return fg
}

func (f *FalcoGateway) openPipe() {
	var err error

	if f.mode == "online" {
		f.pipeFile, err = os.OpenFile(f.inputFileName, os.O_RDONLY, os.ModeNamedPipe)
	}

	if f.mode == "offline" {
		f.pipeFile, err = os.Open(f.inputFileName)
	}

	if err != nil {
		log.Fatal("Open named pipe file error")
	}

	f.pipeReader = bufio.NewReader(f.pipeFile)
	if f.pipeReader == nil {
		log.Fatal("bufio reader opening error")
	}
}

func (f *FalcoGateway) getLine() string {
	line, err := f.pipeReader.ReadBytes('\n')

	if err != nil {
		log.Fatal("error, pipe file broken")
	}

	fmt.Println(string(line))

	return string(line)
}

func (f *FalcoGateway) sendSigRcvSummary() {
	err := syscall.Kill(f.falcoPid, f.sigRcvSummary)
	if err != nil {
		log.Fatalf("Unable to send signal %d to Falco", rcvStatsSignal)
	}
}

func (f *FalcoGateway) sendSigRcvRulesNames() {
	err := syscall.Kill(f.falcoPid, f.sigRcvRulesNames)
	if err != nil {
		log.Fatalf("Unable to send signal %d to Falco", rcvRulesNamesSignal)
	}
}

func (f *FalcoGateway) sendSigFlushData() {
	err := syscall.Kill(f.falcoPid, f.sigFlushData)
	if err != nil {
		log.Fatalf("Unable to send signal %d to Falco", flushDataSignal)
	}
}
