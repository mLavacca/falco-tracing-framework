package main

import (
	"bufio"
	"log"
	"os"
	"syscall"
)

const rcvStatsSignal = 34
const flushDataSignal = 35
const rcvRulesNamesSignal = 36

type FalcoInterface struct {
	falcoPid int
	pipeName string
	pipeFile *os.File

	sigRcvSummary    syscall.Signal
	sigRcvRulesNames syscall.Signal
	sigFlushData     syscall.Signal

	pipeReader *bufio.Reader
}

func NewFalcoInterface(falcoPid int, pipeName string) *FalcoInterface {
	f := new(FalcoInterface)

	f.falcoPid = falcoPid
	f.pipeName = pipeName

	f.sigRcvSummary = syscall.Signal(rcvStatsSignal)
	f.sigRcvRulesNames = syscall.Signal(rcvRulesNamesSignal)
	f.sigFlushData = syscall.Signal(flushDataSignal)

	return f
}

func (f *FalcoInterface) OpenPipe() {
	var err error

	f.pipeFile, err = os.OpenFile(f.pipeName, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		log.Fatal("Open named pipe file error")
	}

	f.pipeReader = bufio.NewReader(f.pipeFile)
	if f.pipeReader == nil {
		log.Fatal("bufio reader opening error")
	}
}

func (f FalcoInterface) getLine() string {
	line, err := f.pipeReader.ReadBytes('\n')

	if err != nil {
		log.Fatal("error, pipe file broken")
	}

	return string(line)
}

func (f FalcoInterface) sendSigRcvSummary() {
	err := syscall.Kill(f.falcoPid, f.sigRcvSummary)
	if err != nil {
		log.Fatalf("Unable to send signal %d to Falco", rcvStatsSignal)
	}
}

func (f FalcoInterface) sendSigRcvRulesNames() {
	err := syscall.Kill(f.falcoPid, f.sigRcvRulesNames)
	if err != nil {
		log.Fatalf("Unable to send signal %d to Falco", rcvRulesNamesSignal)
	}
}

func (f FalcoInterface) sendSigFlushData() {
	err := syscall.Kill(f.falcoPid, f.sigFlushData)
	if err != nil {
		log.Fatalf("Unable to send signal %d to Falco", flushDataSignal)
	}
}
