package falco_test

import (
	"log"
	"os"
	"path"
)

var testFunctions = map[int]interface{}{
	0: writeBelowRoot,
	1: writeBelowEtc,
	2: modifyShellConfigurationFile,
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
	openFile(path, os.O_RDWR|os.O_CREATE)
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
	openFile(path, os.O_RDWR|os.O_CREATE)
}

func writeBelowEtcRollback() {
	path := "/etc/falco_tester_file"
	deleteFile(path)
}

/*
 * id = 2
 */
func modifyShellConfigurationFile() {
	username, err := getUsername()
	if err != nil {
		log.Fatalln(err)
	}

	path := path.Join("/home", username, "/.bashrc")
	openFile(path, os.O_RDWR)
}
