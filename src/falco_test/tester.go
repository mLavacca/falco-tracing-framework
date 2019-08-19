package falco_test

import (
	"configuration"
)

const (
	N = 10000000
	K = 10
)

type TesterFunction struct {
	function interface{}
	counter  int
}

type Tester struct {
	functionList [][]TesterFunction
	rollbackList []int
	ratio        int
	limit        int
}

func NewTester(conf configuration.TracerConfigurations) (*Tester, int) {
	t := new(Tester)
	d := 0

	profile := conf.Record.BreakingProfile

	for _, p := range conf.BreakingProfiles {
		if p.Name == profile {

			t.ratio = p.Ratio
			t.limit = p.Limit
			d = p.Duration

			i := 0
			for _, s := range p.Sequence {
				t.functionList = append(t.functionList, []TesterFunction{})

				for _, v := range s {
					t.functionList[i] = append(t.functionList[i], TesterFunction{
						function: testFunctions[v],
					})
				}
				i++
			}

			for _, r := range p.RollbackSequence {
				t.rollbackList = append(t.rollbackList, r)
			}
			break
		}
	}

	if len(t.functionList) == 0 {
		return nil, d
	} else {
		return t, 0
	}

}

func (t *Tester) RunAllTests() {
	for _, ts := range t.functionList {
		t.runTestSequence(ts)
	}
}

func (t *Tester) RunAllRollbacks() {
	for _, r := range t.rollbackList {
		testRollbacks[r].(func())()
	}
}

func (t *Tester) runTestSequence(testSequence []TesterFunction) {
	n := N
	k := K

	for n > (k * t.ratio) {
		r := (k - (k / n)) / t.ratio

		for i := 0; i < r; i++ {
			for j := 0; j < n; j++ {
			}

			for p := 0; p < len(testSequence); p++ {
				testSequence[p].function.(func())()
				testSequence[p].counter++

				if testSequence[p].counter >= t.limit && p == len(testSequence)-1 {
					return
				}
			}
		}

		n = (n / (K / t.ratio))
		k = (k * K / t.ratio)
	}

}
