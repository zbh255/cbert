package ioc

import (
	"github.com/zbh255/cbert/common/log"
	"io"
	"os"
)

/*
	The ioc library Provide some better abstractions for code that needs to use some basic libraries
	And there is no need to worry about dependent registration
*/
var (
	loggers = func() map[string]log.Logger {
		return make(map[string]log.Logger)
	}()
)

func GetStdLogger() log.Logger {
	return loggers["stdOut"]
}

func GetAccessLogger() log.Logger {
	return loggers["accessLog"]
}

func GetErrorLogger() log.Logger {
	return loggers["errorLog"]
}

func RegisterAccessLogger(writer io.Writer) {
	loggers["accessLog"] = log.NewLogger(writer,log.DEBUG)
}

func RegisterErrorLogger(writer io.Writer) {
	loggers["errorLog"] = log.NewLogger(writer,log.ERROR)
}

func init() {
	loggers["stdOut"] = log.NewLogger(os.Stdout,log.ERROR)
}
