package main

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
	ratio        int
	limit        int
}

func NewTester(conf TracerConfigurations) *Tester {
	t := new(Tester)

	profile := conf.Record.BreakingProfile

	for _, p := range conf.BreakingProfiles {
		if p.Name == profile {

			t.ratio = p.Ratio
			t.limit = p.Limit

			i := 0
			for _, s := range p.Sequence {
				t.functionList = append(t.functionList, []TesterFunction{})

				for _, v := range s {
					t.functionList[i] = append(t.functionList[i], TesterFunction{
						function: functionsSlice[v],
					})
				}
				i++
			}
			break
		}
	}

	return t
}

func (t *Tester) runAllTests() {
	for _, ts := range t.functionList {
		t.runTestSequence(ts)
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

			index := i % len(testSequence)
			testSequence[index].function.(func())()
			testSequence[index].counter++

			if testSequence[index].counter >= t.limit {
				return
			}
		}

		n = (n / (K / t.ratio))
		k = (k * K / t.ratio)
	}

}
