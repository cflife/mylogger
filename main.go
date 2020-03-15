package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"project1/conf"
	"project1/etcd"
	"project1/kafka"
	"project1/taillog"
	"project1/utils"
	"sync"
	"time"
)

var (
	cfg = new(conf.AppConf)
)

func main() {

	err := ini.MapTo(cfg, "./conf/config.ini")
	if err != nil {
		fmt.Printf("load ini faild,err:%v\n", err)
	}
	//1.初始化kafka连接
	err = kafka.Init([]string{cfg.KafkaConf.Address}, cfg.KafkaConf.ChanMaxSize)
	if err != nil {
		fmt.Printf("init kafka failed,err:%v\n", err)
		return
	}
	fmt.Println("kafka init success")

	//2.初始化etcd
	etcd.Init(cfg.EtcdConf.Address, time.Duration(cfg.EtcdConf.TimeOut)*time.Second)
	if err != nil {
		fmt.Printf("init etcd failed,error:%v\n", err)
		return
	}
	fmt.Println("init etcd succcess")

	ipStr, err := utils.GetOutboundIP()
	if err != nil {
		panic(err)
	}
	etcdConfKey := fmt.Sprintf(cfg.EtcdConf.Key, ipStr)
	//从etcd中获取日志收集项的配置信息
	logEntryConf, err := etcd.GetConf(etcdConfKey)
	if err != nil {
		fmt.Printf("etcd.GetConf failed,err:%v\n", err)
		return
	}
	fmt.Printf("get conf from etcd success, %v\n", logEntryConf)

	for index, value := range logEntryConf {
		fmt.Printf("index:%v value:%v\n", index, value)
	}

	//3.收集日志发往kafka
	taillog.Init(logEntryConf)
	// 因为NewConfChan访问了tskMgr的newConfChan, 这个channel是在taillog.Init(logEntryConf) 执行的初始化
	newConfChan := taillog.NewConfChan() // 从taillog包中获取对外暴露的通道
	var wg sync.WaitGroup
	wg.Add(1)
	go etcd.WatchConf(etcdConfKey, newConfChan) // 哨兵发现最新的配置信息会通知上面的那个通道
	wg.Wait()
}
