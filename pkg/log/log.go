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

/*===========================================================================*/
func SetLoggerLevel(s string) {
	lv, err := logrus.ParseLevel(s)
	if err != nil {
		panic(err)
	}
	globalLogger.SetLevel(lv)
}
