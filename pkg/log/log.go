package log

import (
	"github.com/sirupsen/logrus"
)

var globalLogger = logrus.New()

func init() {
	logrus.SetReportCaller(true)
}

func Get(name string) *logrus.Entry {
	return globalLogger.WithField("name", name)
}
