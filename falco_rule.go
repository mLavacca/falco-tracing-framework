package main

import (
	"strconv"
	"strings"
)

type FalcoRule struct {
	id   int
	name string
	tag  int
}

func NewRule(line string) *FalcoRule {
	line = strings.Replace(line, "\n", "", 1)

	tracerLine := strings.Split(line, "-")

	id, err := strconv.Atoi(tracerLine[2])
	if err != nil {
		return nil
	}
	name := tracerLine[1]

	r := new(FalcoRule)
	r.id = id
	r.name = name

	return r
}
