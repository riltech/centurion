package logger

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func LogError(err error) {
	if err, ok := err.(stackTracer); ok {
		logrus.Errorf("%+v", err)
		return
	}
	logrus.Errorf("%+v", errors.WithStack(err))
}
