// 这里主要是kafak的相关操作，包括了kafka的初始化，以及发送消息的操作
package main

import (
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

var (
	client sarama.SyncProducer
	kafkaSender *KafkaSender
)

type KafkaSender struct {
	client sarama.SyncProducer
	lineChan chan string
}

// 初始化kafka
func NewKafkaSender(kafkaAddr string)(kafka *KafkaSender,err error){
	kafka = &KafkaSender{
		lineChan:make(chan string,100000),
	}
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	client,err := sarama.NewSyncProducer([]string{kafkaAddr},config)
	if err != nil{
		logs.Error("init kafka client failed,err:%v\n",err)
		return
	}
	kafka.client = client
	for i:=0;i<appConfig.KafkaThreadNum;i++{
		// 根据配置文件循环开启线程去发消息到kafka
		go kafka.sendToKafka()
	}
	return
}

func initKafka()(err error){
	kafkaSender,err = NewKafkaSender(appConfig.kafkaAddr)
	return
}

func (k *KafkaSender) sendToKafka(){
	//从channel中读取日志内容放到kafka消息队列中
	for v := range k.lineChan{
		msg := &sarama.ProducerMessage{}
		msg.Topic = "nginx_log"
		msg.Value = sarama.StringEncoder(v)
		_,_,err := k.client.SendMessage(msg)
		if err != nil{
			logs.Error("send message to kafka failed,err:%v",err)
		}
	}
}

func (k *KafkaSender) addMessage(line string)(err error){
	//我们通过tailf读取的日志文件内容先放到channel里面
	k.lineChan <- line
	return
}

