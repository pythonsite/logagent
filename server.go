package main

import (
	"github.com/hpcloud/tail"
	"fmt"
	"sync"
	"github.com/astaxie/beego/logs"
	"strings"
)

type TailMgr struct {
	//因为我们的agent可能是读取多个日志文件，这里通过存储为一个map
	tailObjMap map[string]*TailObj
	lock sync.Mutex
}

type TailObj struct {
	//这里是每个读取日志文件的对象
	tail *tail.Tail
	offset int64  //记录当前位置
	filename string
}

var tailMgr *TailMgr
var waitGroup sync.WaitGroup

func NewTailMgr()(*TailMgr){
	tailMgr =  &TailMgr{
		tailObjMap:make(map[string]*TailObj,16),
	}
	return tailMgr
}

func (t *TailMgr) AddLogFile(filename string)(err error){
	t.lock.Lock()
	defer t.lock.Unlock()
	_,ok := t.tailObjMap[filename]
	if ok{
		err = fmt.Errorf("duplicate filename:%s\n",filename)
		return
	}
	tail,err := tail.TailFile(filename,tail.Config{
		ReOpen:true,
		Follow:true,
		Location:&tail.SeekInfo{Offset:0,Whence:2},
		MustExist:false,
		Poll:true,
	})

	tailobj := &TailObj{
		filename:filename,
		offset:0,
		tail:tail,
	}
	t.tailObjMap[filename] = tailobj
	return
}

func (t *TailMgr) Process(){
	//开启线程去读日志文件
	for _, tailObj := range t.tailObjMap{
		waitGroup.Add(1)
		go tailObj.readLog()
	}
}

func (t *TailObj) readLog(){
	//读取每行日志内容
	for line := range t.tail.Lines{
		if line.Err != nil {
			logs.Error("read line failed,err:%v",line.Err)
			continue
		}
		str := strings.TrimSpace(line.Text)
		if len(str)==0 || str[0] == '\n'{
			continue
		}

		kafkaSender.addMessage(line.Text)
	}
	waitGroup.Done()
}


func RunServer(){
	tailMgr = NewTailMgr()
	// 这一部分是要调用tailf读日志文件推送到kafka中
	for _, filename := range appConfig.LogFiles{
		err := tailMgr.AddLogFile(filename)
		if err != nil{
			logs.Error("add log file failed,err:%v",err)
			continue
		}

	}
	tailMgr.Process()
	waitGroup.Wait()
}






