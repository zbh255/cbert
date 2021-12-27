package main

import (
	"flag"
	"github.com/zbh255/cbert/connect"
	"github.com/zbh255/cbert/ioc"
	"net"
	"os"
)

func main() {
	// register
	configFilePath := flag.String("projectConfig", "./project_config.json", "主配置文件")
	flag.Parse()
	err := ioc.RegisterProjectConfig(*configFilePath)
	if err != nil {
		panic(err)
	}
	projectConfig := ioc.GetProjectConfig()
	//
	if err = ioc.RegisterUserConfig(projectConfig.Source.UserConfig); err != nil {
		panic(err)
	}
	accessLogWriter, err := os.OpenFile(projectConfig.Log.AccessLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	defer accessLogWriter.Close()
	errorLogWriter, err := os.OpenFile(projectConfig.Log.ErrorLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	defer errorLogWriter.Close()
	ioc.RegisterAccessLogger(accessLogWriter)
	ioc.RegisterErrorLogger(errorLogWriter)
	// to live recover error
	defer func() {
		if err := recover(); err != nil {
			errLog := ioc.GetErrorLogger()
			errLog.Error(err.(error).Error())
			panic(err)
		}
	}()
	// handler user source
	handlerSource()
	// register command handler
	handleGenerator()
	// handler user add source request
	handlerAddSource()
	mainStart()
}

func mainStart() {
	projectConfig := ioc.GetProjectConfig()
	listener, err := net.Listen("tcp",projectConfig.Connection.Addr)
	if err != nil {
		panic(err)
	}
	server := connect.NewConnection(listener)
	_ = server.Start()
}