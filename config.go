package main

import (
	"go_dev/13/config"
	"strings"
	"fmt"
	"github.com/astaxie/beego/logs"
)

type AppConfig struct {
	LogPath string
	LogLevel string
	kafkaAddr string
	KafkaThreadNum int
	LogFiles []string
	etcdAddr []string
	etcdWatchKeyFmt string
	//单位为毫秒
	etcdTimeout int
}

var appConfig = &AppConfig{}

func initConfig(filename string)(err error){
	conf,err := config.NewConfig(filename)
	if err != nil{
		return
	}
	logs.Info(conf.GetString("etcd_watch_key"))
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


	etcdAddr, err := conf.GetString("etcd_addr")
	if err != nil || len(etcdAddr) == 0{
		return
	}

	arr := strings.Split(etcdAddr, ",")
	for _, v := range arr {
		str := strings.TrimSpace(v)
		if len(str) == 0 {
			continue
		}
		appConfig.etcdAddr = append(appConfig.etcdAddr, str)
	}

	etcdKey, err := conf.GetString("etcd_watch_key")
	if err != nil || len(etcdKey) == 0{
		logs.Warn("get etcd watch key failed, err:%v", err)
		return
	}

	appConfig.etcdTimeout = conf.GetIntDefault("etcd_timeout", 1500)
	appConfig.KafkaThreadNum = conf.GetIntDefault("kafka_thread_num",8)

	appConfig.etcdWatchKeyFmt = etcdKey
	appConfig.kafkaAddr = kafkaAddr
	appConfig.LogLevel = logLevel
	appConfig.LogPath = logPath

	fmt.Printf("load config succ, data:%v\n", appConfig)
	return
}
