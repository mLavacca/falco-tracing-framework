package falco_test

import (
	"log"
	"os"
	"path"
)

var testFunctions = map[int]interface{}{
	3:  modifyShellConfigurationFile,
	6:  updatePackageRepository,
	7:  writeBelowBinaryDir,
	8:  writeBelowMonitoredDir,
	10: writeBelowEtc,
	11: writeBelowRoot,
	12: readSensitiveFilesAfterStartup,
	13: readSensitiveFileUntrusted,
}

var testRollbacks = map[int]interface{}{
	7:  writeBelowBinaryDirRollback,
	9:  writeBelowMonitoredDirRollback,
	10: writeBelowEtcRollback,
	11: writeBelowRootRollback,
}

func writeBelowRoot() {
	filePath := "/falco_tester_file"
	openFile(filePath, os.O_RDWR|os.O_CREATE)
}

func writeBelowRootRollback() {
	filePath := "/falco_tester_file"
	deleteFile(filePath)
}

func writeBelowEtc() {
	filePath := "/etc/falco_tester_file"
	openFile(filePath, os.O_RDWR|os.O_CREATE)
}

func writeBelowEtcRollback() {
	filePath := "/etc/falco_tester_file"
	deleteFile(filePath)
}

func modifyShellConfigurationFile() {
	username, err := getUsername()
	if err != nil {
		log.Fatalln(err)
	}

	filePath := path.Join("/home", username, "/.bashrc")
	openFile(filePath, os.O_RDWR)
}

func updatePackageRepository() {
	filePath := "/etc/apt/sources.list"
	openFile(filePath, os.O_RDWR)
}

func writeBelowBinaryDir() {
	filePath := "/usr/sbin/falco_tester_file"
	openFile(filePath, os.O_RDWR|os.O_CREATE)
}

func writeBelowBinaryDirRollback() {
	filePath := "/usr/sbin/falco_tester_file"
	deleteFile(filePath)
}

func writeBelowMonitoredDir() {
	filePath := "/usr/local/bin/falco_tester_file"
	openFile(filePath, os.O_RDWR|os.O_CREATE)
}

func writeBelowMonitoredDirRollback() {
	filePath := "/usr/local/bin/falco_tester_file"
	deleteFile(filePath)
}

func readSensitiveFilesAfterStartup() {

}

func readSensitiveFileUntrusted() {
	filePath := "/etc/pam.d/login"
	openFile(filePath, os.O_RDONLY)
}
