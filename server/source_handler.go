package main

import (
	"flag"
	"fmt"
	"github.com/zbh255/cbert/ioc"
	"io/ioutil"
	"os"
	"strings"
)

var (
	fileName = flag.String("fileName", "", "要添加的文件名")
	uuidStr  = flag.String("uuid", "", "操作指定uuid下的资源")
)

// Prepare corresponding resources for users
func handlerSource() {
	projectConfig := ioc.GetProjectConfig()
	dir, err := ioutil.ReadDir(projectConfig.Source.Dir)
	if err != nil {
		panic(err)
	}
	hashTable := make(map[string]int, len(dir))
	for k, v := range dir {
		hashTable[v.Name()] = k
	}
	userConfig := ioc.GetUserConfig()
	for _, user := range userConfig.User {
		k, ok := hashTable[user]
		if ok && dir[k].IsDir() {
			continue
		}
		if err = os.Mkdir(projectConfig.Source.Dir+"/"+user, 0755); err != nil {
			panic(err)
		}
	}
}

func handlerAddSource() {
	if *fileName == "" || *uuidStr == "" {
		return
	}
	projectConfig := ioc.GetProjectConfig()
	userConfig := ioc.GetUserConfig()
	buffer := make([]string, 0, 1)
	for _, v := range userConfig.User {
		if strings.EqualFold(*uuidStr, v[:len(*uuidStr)]) {
			buffer = append(buffer, v)
		}
	}
	if len(buffer) > 1 {
		fmt.Printf("Which is the real choice?\n")
		for _, v := range buffer {
			fmt.Println(v)
		}
		return
	}
	filePath := projectConfig.Source.Dir + "/" + buffer[0] + "/" + *fileName
	errLog := ioc.GetErrorLogger()
	accessLog := ioc.GetAccessLogger()
	f, err := os.Create(filePath)
	if err != nil {
		errLog.Error(err.Error())
		return
	} else {
		f.Close()
		accessLog.Info(fmt.Sprintf("%s create file: %s successfully", buffer[0], *fileName))
	}
}
