package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"encoding/json"
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

	err = initKafka()
	if err != nil{
		return
	}
	RunServer()
}
