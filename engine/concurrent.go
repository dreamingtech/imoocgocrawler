package engine

import (
	"log"
)

type ConcurrentEngine struct {
	// 3. Engine.Scheduler
	Scheduler   iScheduler
	WorkerCount int
}

// 2. iScheduler 接口, 实现 Submit 方法
type iScheduler interface {
	Submit(Request)
	// ConfigureWorkerChan 的作用
	// Scheduler 中需要有一个 channel, 用来保存 request,
	// 此队列即为 Run 中创建的 in channel
	// 为了能把外部创建的队列添加到 Scheduler 中, 要添加一个方法,
	// 通过此方法把外部创建的 channel 赋值给 Scheduler
	ConfigureWorkerChan(chan Request)
}

func (engine *ConcurrentEngine) Run(seeds ...Request) {

	// 5. 所有的 Worker 共用一个输入队列, 一个输出队列
	// Worker 从 in 中取 Request, 把解析到的数据保存到 out 中
	in := make(chan Request)
	out := make(chan ParseResult)

	// 把创建的 in channel 保存到 Scheduler 中
	engine.Scheduler.ConfigureWorkerChan(in)

	// 4. 创建 Worker
	for i := 0; i < engine.WorkerCount; i++ {
		// 创建一个 Worker
		createWorker(in, out)
	}

	// 1. 调用 Scheduler 中的方法, 将接收到的所有 Request 添加到队列中
	for _, r := range seeds {
		engine.Scheduler.Submit(r)
	}

	// 从 Out Channel 中取数据
	for {
		result := <-out
		for _, item := range result.Items {
			log.Printf("Got item: %v", item)
		}

		// 把 result 中的请求提交给 Scheduler
		for _, request := range result.Requests {
			engine.Scheduler.Submit(request)
		}
	}
}

// todo 是不是叫 doWork 更好点
func createWorker(in chan Request, out chan ParseResult) {
	go func() {
		for {
			// 从 in channel 中取出 request, 交给 worker 处理, 并把 worker 处理的结果送入到 out 中
			request := <-in
			// 可以使用 simple engine 中的 worker 来发送请求, 处理响应
			result, err := worker(request)
			if err != nil {
				continue
			}
			out <- result
		}
	}()
}
