package main

import (
	"fmt"
	"go.uber.org/zap"
	"report/common/config"
	"report/common/db/mgClient"
	"report/common/log"
	"report/handler/antiHandler"
	"time"
)

func init() {
	fmt.Println("读取配置文件...")
	config.Init()

	//初始化日志
	fmt.Println("初始化日志...")
	log.GetLogger().Info("init zap logger success!")
}

func main() {
	//连接mongo数据库
	fmt.Println("连接mongoDB...")
	err := mgClient.Init()
	if err != nil {
		log.GetLogger().Error("init mongoDB client failed!", zap.Error(err))
		panic(err)
	}

	fmt.Println("程序启动完成[OK]")

	start()
}

func start() {
	//每天凌晨1点开始执行
	trigTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 1, 0, 0, 0, time.Local)

	//启动时间如果已经过了执行时间，那么第二天再执行
	if time.Now().After(trigTime) {
		trigTime = trigTime.AddDate(0, 0, 1)
	}

	for {
		now := time.Now()
		if now.After(trigTime) && now.Hour() == trigTime.Hour() {
			//执行任务
			log.GetLogger().Info("开始执行报表任务start")

			//开始执行anti-ban的报表统计
			antiHandler.Start()

			//重置时间
			log.GetLogger().Info("任务执行结束over")
			trigTime = trigTime.AddDate(0, 0, 1)
		}

		//每分钟检查一次时间
		time.Sleep(1 * time.Minute)
	}
}
