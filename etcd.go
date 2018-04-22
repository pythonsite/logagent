package main

import (
	"github.com/coreos/etcd/clientv3"
	"time"
	"github.com/astaxie/beego/logs"
	"context"
	"fmt"
)

var Client *clientv3.Client
var logConfChan chan string


// 初始化etcd
func initEtcd(addr []string,keyfmt string,timeout time.Duration)(err error){

	var keys []string
	for _,ip := range ipArrays{
		//keyfmt = /logagent/%s/log_config
		keys = append(keys,fmt.Sprintf(keyfmt,ip))
	}

	logConfChan = make(chan string,10)
	logs.Debug("etcd watch key:%v timeout:%v", keys, timeout)

	Client,err = clientv3.New(clientv3.Config{
		Endpoints:addr,
		DialTimeout: timeout,
	})
	if err != nil{
		logs.Error("connect failed,err:%v",err)
		return
	}
	logs.Debug("init etcd success")
	waitGroup.Add(1)
	for _, key := range keys{
		ctx,cancel := context.WithTimeout(context.Background(),2*time.Second)
		// 从etcd中获取要收集日志的信息
		resp,err := Client.Get(ctx,key)
		cancel()
		if err != nil {
			logs.Warn("get key %s failed,err:%v",key,err)
			continue
		}

		for _, ev := range resp.Kvs{
			logs.Debug("%q : %q\n",  ev.Key, ev.Value)
			logConfChan <- string(ev.Value)
		}
	}
	go WatchEtcd(keys)
	return
}

func WatchEtcd(keys []string){
	// 这里用于检测当需要收集的日志信息更改时及时更新
	var watchChans []clientv3.WatchChan
	for _,key := range keys{
		rch := Client.Watch(context.Background(),key)
		watchChans = append(watchChans,rch)
	}

	for {
		for _,watchC := range watchChans{
			select{
			case wresp := <-watchC:
				for _,ev:= range wresp.Events{
					logs.Debug("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
					logConfChan <- string(ev.Kv.Value)
				}
			default:

			}
		}
		time.Sleep(time.Second)
	}
	waitGroup.Done()
}

func GetLogConf()chan string{
	return logConfChan
}