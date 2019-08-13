package stats_getter

import (
	"bufio"
	"log"
	"os"
	"syscall"
)

const (
	rcvStatsSignal      = 34
	flushDataSignal     = 35
	rcvRulesNamesSignal = 36
)

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

func newFalcoGateway(mode string) *FalcoGateway {

	fg := new(FalcoGateway)
	fg.mode = mode

	if mode == "online" {
		fg.configureSignals()
	}

	fg.openPipe()

	return fg
}

func (f *FalcoGateway) openPipe() {
	var err error

	if f.mode == "online" {
		f.inputFileName = "/tmp/falco_tracer_pipe"
		f.pipeFile, err = os.OpenFile(f.inputFileName, os.O_RDWR, os.ModeNamedPipe)
	}

	if f.mode == "offline" {
		f.inputFileName = "/tmp/falco_tracer_file"
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

func (fg *FalcoGateway) closePipe() {
	err := fg.pipeFile.Close()
	if err != nil {
		log.Fatalln("error during file closing", err)
	}

	err = os.Remove(fg.inputFileName)
	if err != nil {
		log.Fatalln("error during file removing", err)
	}
}

func (f *FalcoGateway) getLine() string {
	line, err := f.pipeReader.ReadBytes('\n')

	if err != nil {
		log.Fatal("error, pipe file broken")
	}

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
