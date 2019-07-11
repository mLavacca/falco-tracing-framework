package main

type RuleStat struct {
	name    string
	tag     int
	counter uint64
	latency uint64
}

func NewRuleStat(name string, tag int, counter uint64, latency uint64) *RuleStat {
	r := new(RuleStat)

	r.name = name
	r.tag = tag
	r.counter = counter
	r.latency = latency

	return r
}

type FuncStat struct {
	name    string
	counter uint64
	latency uint64
}

func NewFuncStat(name string, counter uint64, latency uint64) *FuncStat {
	f := new(FuncStat)

	f.name = name
	f.counter = counter
	f.latency = latency

	return f
}

type CounterStat struct {
	name    string
	counter uint64
}

func NewCounterStat(name string, counter uint64) *CounterStat {
	c := new(CounterStat)

	c.name = name
	c.counter = counter

	return c
}

type Falcostats struct {
	unbrokenRules []RuleStat
	brokenRules   []RuleStat
	funcStats     []FuncStat
	counterStats  []CounterStat
}
