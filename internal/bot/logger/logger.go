package logger

import (
	"github.com/sirupsen/logrus"
)

func New() *logrus.Logger {
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)

	return l
}
