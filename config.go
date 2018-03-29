package main

import (
	"go_dev/13/config"
	"strings"
	"fmt"
)

type AppConfig struct {
	LogPath string
	LogLevel string
	kafkaAddr string
	KafkaThreadNum int
	LogFiles []string
}

var appConfig = &AppConfig{}

func initConfig(filename string)(err error){
	conf,err := config.NewConfig(filename)
	if err != nil{
		return
	}
	logPath,err := conf.GetString("log_path")
	if err != nil || len(logPath) == 0 {
		return
	}
	logLevel,err := conf.GetString("log_level")
	if err != nil || len(logLevel) == 0 {
		return
	}
	kafkaAddr,err:= conf.GetString("kafka_addr")
	if err != nil || len(kafkaAddr) == 0 {
		return
	}
	logFiles,err := conf.GetString("log_files")
	if err != nil || len(logFiles) == 0 {
		return
	}

	appConfig.KafkaThreadNum = conf.GetIntDefault("kafka_thread_num",8)

	arr := strings.Split(logFiles,",")
	for _,v := range arr{
		str := strings.TrimSpace(v)
		if len(str) == 0{
			continue
		}
		appConfig.LogFiles = append(appConfig.LogFiles,str)
	}
	appConfig.kafkaAddr = kafkaAddr
	appConfig.LogLevel = logLevel
	appConfig.LogPath = logPath

	fmt.Printf("load conf success data:%v\n",appConfig)
	return
}
