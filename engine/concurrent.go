package engine

import (
	"log"
)

type ConcurrentEngine struct {
	// 3. Engine.Scheduler
	Scheduler   iScheduler
	WorkerCount int
	// 定义 item channel, item 可以是任意类型的数据
	ItemChan chan Item
}

// 2. iScheduler 接口, 实现 Submit 方法
type iScheduler interface {
	Submit(Request)
	// ConfigureWorkerChan 的作用
	// Scheduler 中需要有一个 channel, 用来保存 request,
	// 此队列即为 Run 中创建的 in channel
	// 为了能把外部创建的队列添加到 Scheduler 中, 要添加一个方法,
	// 通过此方法把外部创建的 channel 赋值给 Scheduler
	// ConfigureWorkerChan(chan Request)

	// GetWorkerChan 统一 simple scheduler 和 queued scheduler 的代码
	// simple scheduler 和 queued scheduler 的主要区别是:
	// simple scheduler 是所有的 worker 共用一个 worker channel,
	// 使用 ConfigureWorkerChan 来设置这个统一的 worker channel.
	// queued scheduler 是每个 worker 有自己的 worker channel
	// 又有一个总的 worker channel 保存所有的 worker channel,
	// 所以 queued scheduler 不需要 ConfigureWorkerChan 来设置 worker channel.
	// 但是使用同一个 workerChan 还是每一个 worker 有自己的 workerChan, 这是属于 scheduler 的功能,
	// 应该 scheduler 中完成 workerChan 的创建和配置, 向 scheduler 提交 request 的时候,
	// 向 scheduler 要一个 workerChan, 由 scheduler 确定使用的 workerChan
	GetWorkerChan() chan Request

	// iReadyNotifier 使用接口组合, 把 ReadyNotifier 接口嵌入到 iScheduler 接口中
	iReadyNotifier
	// WorkerReady(chan Request)
	Run()
}

// ReadyNotifier 的作用: iScheduler 中的 4 个方法, 比较多, 可以把 WorkerReady 方法提取成一个新接口
// 此时接口的实现者不需要修改, 就已经实现了这两个接口中的所有方法, 也就实现了这两个接口
type iReadyNotifier interface {
	WorkerReady(chan Request)
}

func (engine *ConcurrentEngine) Run(seeds ...Request) {

	// 5. 所有的 Worker 共用一个输入队列, 一个输出队列
	// Worker 从 in 中取 Request, 把解析到的数据保存到 out 中
	// 因为所有的 Request 队列都已经在 Scheduler 中实现了, 这里就不用再生成 request channel 了
	// in := make(chan Request)

	out := make(chan ParseResult)

	// 把创建的 in channel 保存到 Scheduler 中
	// engine.Scheduler.ConfigureWorkerChan(in)
	// 此时就不再是 runQueuedConcurrentEngine 了, 而是运行 Scheduler
	// Scheduler 会创建 workerChan
	engine.Scheduler.Run()

	// 4. 创建 Worker
	for i := 0; i < engine.WorkerCount; i++ {
		// 创建 worker 时, 向 Scheduler "要" 一个 workerChan
		workerChan := engine.Scheduler.GetWorkerChan()
		// 创建一个 Worker
		// 因为 Request 队列是在 Scheduler 中的, 要把 Scheduler 传递给 createWorker
		createWorker(engine.Scheduler, workerChan, out)
	}

	// 1. 调用 Scheduler 中的方法, 将接收到的所有 Request 添加到队列中
	for _, r := range seeds {
		// 因为去重首先要保存已经 "见到" 过的 url, 所以也要对 seeds 进行去重
		if isDuplicate(r.Url) {
			log.Printf("Duplicate request: %s", r.Url)
			continue
		}
		engine.Scheduler.Submit(r)
	}

	// 循环等待出现的原因及解决方法
	/*
		Engine 通过函数调用把 request 提交给 Scheculer, `engine.Scheduler.Submit(r)`,
		Scheduler 通过向 workerChan 中发送 request 来实现任务的分发, `s.workerChan <- request`
		Scheduler 向 workerChan 中发送数据成功的前提是 有一个空闲的 worker 在等待从 workerChan 中收取 request, `request := <-in`
		worker 等待从 workerChan 中收取 request 的前提是 `把上一件事情做完`, 即把上一次请求中解析到的 request 和 item 发送给 engine,
		engine 再调用 Scheduler 向 workerChan 中发送 request, `out <- result`,
		但向 workerChan 中发送请求成功的前提是要有一个空闲的 Worker, 这样就形成了一个循环等待
		所以前 10 个请求发送出去之后, 程序就会陷入到循环等待中, 也就是卡死了
		只需要使用 goroutine 把向 workerChan 中提交 request 的操作变成异步的, 就解决了以上问题
		此时, engine 再从 out 中取数据, `result := <-out`, 提交给 Scheduler 时, `engine.Scheduler.Submit(request)`,
		因为 Submit 变成了协程, 即异步的方式执行, 就不会再出现循环等待的问题了
	*/

	// 从 Out Channel 中取数据
	for {
		result := <-out
		for _, item := range result.Items {
			// 直接把 item  中
			go func(i Item) {
				engine.ItemChan <- i
			}(item)
		}

		// 把 result 中的请求提交给 Scheduler
		for _, request := range result.Requests {
			if isDuplicate(request.Url) {
				log.Printf("Duplicate request: %s", request.Url)
				continue
			}
			engine.Scheduler.Submit(request)
		}
	}
}

// todo 是不是叫 doWork 更好点
func createWorker(ready iReadyNotifier, in chan Request, out chan ParseResult) {
	// 每个 worker 都创建一个自己的 channel
	// workerChan 是用来对 Request 进行分发调度的, 是 Scheduler 的功能, 所以不应该在 engine 中创建
	// 不再是创建 worker 的时候创建自己的 channel, 而是在 scheduler 中创建
	// in := make(chan Request)

	go func() {
		for {
			// tell scheduler i'm ready
			// 调用 WorkerReady 时, 把 workerChan 传递过去, 再把 workerChan 保存到 Scheduler 的 workerQ 中
			// 如果有 request 发送给了 workerChan, 就会继续执行下面的操作
			ready.WorkerReady(in)

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

// 要想去重，就需要把已经访问过的 url 保存起来
var visitedUrls = make(map[string]bool)

// isDuplicate 用来判断 url 是否已经访问过
func isDuplicate(url string) bool {
	if visitedUrls[url] {
		return true
	}
	visitedUrls[url] = true
	return false
}
