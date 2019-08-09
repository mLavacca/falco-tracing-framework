package falco_test

var testFunctions = map[int]interface{}{
	0: writeBelowRoot,
	1: writeBelowEtc,
}

var testRollbacks = map[int]interface{}{
	0: writeBelowRootRollback,
	1: writeBelowEtcRollback,
}

/*
 * id = 0
 */
func writeBelowRoot() {
	path := "/falco_tester_file"
	writeFile(path)
}

func writeBelowRootRollback() {
	path := "/falco_tester_file"
	deleteFile(path)
}

/*
 * id = 1
 */
func writeBelowEtc() {
	path := "/etc/falco_tester_file"
	writeFile(path)
}

func writeBelowEtcRollback() {
	path := "/etc/falco_tester_file"
	deleteFile(path)
}
