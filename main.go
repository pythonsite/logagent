package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"encoding/json"
	"time"
)

func getLevel(level string) int{
	switch level{
	case "debug":
		return logs.LevelDebug
	case "trace":
		return logs.LevelTrace
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "error":
		return logs.LevelError
	default:
		return logs.LevelDebug
	}
}

func initLog()(err error){
	//初始化日志库
	config := make(map[string]interface{})
	config["filename"] = "./logs/logcollect.log"
	config["level"] = getLevel(appConfig.LogLevel)
	configStr,err := json.Marshal(config)
	if err != nil{
		fmt.Println("mashal failed,err:",err)
		return
	}
	logs.SetLogger(logs.AdapterFile,string(configStr))
	return
}


func main() {
	err := initConfig("./conf/app.conf")
	if err != nil{
		panic(fmt.Sprintf("init config failed,err:%v\n",err))
	}
	err = initLog()
	if err != nil{
		return
	}
	logs.Debug("init success")
	ipArrays, err = getLocalIP()
	logs.Info(ipArrays)
	if err != nil {
		logs.Error("get local ip failed, err:%v", err)
		return
	}

	logs.Debug("get local ip succ, ips:%v", ipArrays)
	err = initKafka()
	if err != nil {
		logs.Error("init kafka faild, err:%v", err)
		return
	}
	//logs.Info(appConfig.etcdAddr,appConfig.etcdWatchKeyFmt,time.Duration(appConfig.etcdTimeout)*time.Microsecond)
	err = initEtcd(appConfig.etcdAddr, appConfig.etcdWatchKeyFmt,
		time.Duration(appConfig.etcdTimeout)*time.Millisecond)
	if err != nil {
		logs.Error("init etcd failed, err:%v", err)
		return
	}
	RunServer()
}
