package main

import (
	"flag"
	"github.com/zbh255/cbert/common/uuid"
	"github.com/zbh255/cbert/ioc"
)

var (
	generator = flag.String("g","","选择指定的生成器:uuid")
)


// console application generator command
func handleGenerator() {
	errLog := ioc.GetAccessLogger()
	accessLog := ioc.GetAccessLogger()
	switch *generator {
	case "uuid","UUID":
		uuidStr,err := uuid.GetCustomUuid()
		if err != nil {
			errLog.Error(err.Error())
			return
		}
		println(uuidStr)
		accessLog.Info("created uuid successfully")
		return
	}
}
