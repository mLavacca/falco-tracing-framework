package main

type Tester struct {
	functionList []struct {
		function interface{}
		counter  int
	}

	ratio int
	limit int
}

func NewTester(conf TracerConfigurations) *Tester {
	t := new(Tester)

	profile := conf.Record.BreakingProfile

	if profile == "flat" {

		return t
	}

	for _, p := range conf.BreakingProfiles {
		if p.Name == profile {
			t.ratio = p.Ratio
			t.limit = p.Limit

			break
		}
	}

	return t
}
