package commands

import "configuration"

func DispatchCommand(cmd string, conf configuration.TracerConfigurations) {
	switch cmd {
	case "record":
		recorder := NewRecorder(conf)
		recorder.StartRecord()
		recorder.Rollback()
	case "offline-report":
		reporter := newOfflineReporter(conf.OfflineReport)
		reporter.report()
	case "online-report":
		reporter := newOnlineReporter(conf.OnlineReport)
		reporter.report()
	}
}
