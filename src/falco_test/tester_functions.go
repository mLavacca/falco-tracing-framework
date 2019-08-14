package falco_test

import (
	"log"
	"os"
	"path"
)

var testFunctions = map[int]interface{}{
	3:  modifyShellConfigurationFile,
	6:  updatePackageRepository,
	10: writeBelowEtc,
	11: writeBelowRoot,
}

var testRollbacks = map[int]interface{}{
	0: writeBelowRootRollback,
	1: writeBelowEtcRollback,
}

func writeBelowRoot() {
	path := "/falco_tester_file"
	openFile(path, os.O_RDWR|os.O_CREATE)
}

func writeBelowRootRollback() {
	path := "/falco_tester_file"
	deleteFile(path)
}

func writeBelowEtc() {
	path := "/etc/falco_tester_file"
	openFile(path, os.O_RDWR|os.O_CREATE)
}

func writeBelowEtcRollback() {
	path := "/etc/falco_tester_file"
	deleteFile(path)
}

func modifyShellConfigurationFile() {
	username, err := getUsername()
	if err != nil {
		log.Fatalln(err)
	}

	path := path.Join("/home", username, "/.bashrc")
	openFile(path, os.O_RDWR)
}

func updatePackageRepository() {
	path := "/etc/apt/sources.list"
	openFile(path, os.O_RDWR)
}
