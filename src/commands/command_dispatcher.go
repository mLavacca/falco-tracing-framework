package commands

import "configuration"

func DispatchCommand(cmd string, conf configuration.TracerConfigurations) {
	switch cmd {
	case "record":
		recorder := newRecorder(conf)
		recorder.startRecord()
		recorder.rollback()
	case "offline-report":
		reporter := newOfflineReporter(conf.OfflineReport)
		reporter.report()
	case "online-report":
		reporter := newOnlineReporter(conf.OnlineReport)
		reporter.report()
	}
}
