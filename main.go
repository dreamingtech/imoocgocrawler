package main

import (
	"github.com/dreamingtech/imoocgocrawler/engine"
	parser "github.com/dreamingtech/imoocgocrawler/parser/zhenai"
	"github.com/dreamingtech/imoocgocrawler/scheduler"
)

func runSimpleEngine() {
	engine.SimpleEngine{}.Run(engine.Request{
		// Url: "http://www.zhenai.com/zhenghun",
		Url:        "http://localhost:8080/mock/www.zhenai.com/zhenghun",
		ParserFunc: parser.ParseCityList,
	})
}

func runConcurrentEngine() {
	// 因为是指针接收者, 必须要定义一个变量, 只有变量才可以取地址
	concurrentEngine := engine.ConcurrentEngine{
		Scheduler: &scheduler.SimpleScheduler{},
		// 设置运行的 worker 数量
		WorkerCount: 10,
	}
	concurrentEngine.Run(engine.Request{
		Url:        "http://localhost:8080/mock/www.zhenai.com/zhenghun",
		ParserFunc: parser.ParseCityList,
	})
}

func runQueuedSchedulerConcurrentEngine() {
	// 因为是指针接收者, 必须要定义一个变量, 只有变量才可以取地址
	concurrentEngine := engine.ConcurrentEngine{
		Scheduler: &scheduler.QueuedScheduler{},
		// 设置运行的 worker 数量
		WorkerCount: 10,
	}
	concurrentEngine.Run(engine.Request{
		Url:        "http://localhost:8080/mock/www.zhenai.com/zhenghun",
		ParserFunc: parser.ParseCityList,
	})
}

func main() {
	// runSimpleEngine()
	// runConcurrentEngine()
	runQueuedSchedulerConcurrentEngine()
}
