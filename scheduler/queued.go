package scheduler

import "github.com/dreamingtech/imoocgocrawler/engine"

type QueuedScheduler struct {
	// 定义 Request 队列和 Worker 队列
	requestChan chan engine.Request
	// todo 定义 worker 类型
	// workerChan 是 worker 类型, 而 worker 对外的接口又是 chan Request
	// 所以 workerChan 的类型是 chan chan Request
	// 每个 worker 会创建不同的 channel, 所有的 worker channel 会放在总的 workerChan 中
	// 在任务分发时, 会把 Request 分发给 workerChan 中的 chan Request
	workerChan chan chan engine.Request
}

// WorkerReady 外界调用此方法通知 Scheduler 有一个 worker ready 了, 可以去接收 Request 了
func (s *QueuedScheduler) WorkerReady(worker chan engine.Request) {
	s.workerChan <- worker
}

func (s *QueuedScheduler) Submit(request engine.Request) {
	s.requestChan <- request
}

func (s *QueuedScheduler) ConfigureWorkerChan(requests chan engine.Request) {
	// TODO implement me
	panic("implement me")
}

func (s *QueuedScheduler) Run() {
	// 定义 workerChan 和 requestChan
	// 因为要改变 s 的内容, 所以要使用指针接收者
	s.workerChan = make(chan chan engine.Request)
	s.requestChan = make(chan engine.Request)

	go func() {
		// 创建两个队列, 把收到的 Request 和 worker 放在这两个队列中进行排队
		var requestQ []engine.Request
		var workerQ []chan engine.Request
		for {
			// 什么时候可以把 Request 发送给 Worker ??
			// 当 requestQ 队列和 workerQ 队列同时有数据时, 即既有 Request 在排队又有 worker 在排队时

			// 但此时不能直接把 Request 发送给 worker, 因为还可能会出现循环等待的问题
			// 而是把所有对 chan 的操作都放在 select 中执行
			// 定义 activeRequest 和 activeWorker
			var activeRequest engine.Request
			var activeWorker chan engine.Request

			if len(requestQ) > 0 && len(workerQ) > 0 {
				// 从队列中取出第一个数据, 但此时还不能直接取出来第一个数据, 只是作了一个标记
				// 当在 select 中真正能够把 request 放入到 worker 中之后, 才从队列中删除数据
				activeRequest = requestQ[0]
				activeWorker = workerQ[0]
			}

			select {
			case request := <-s.requestChan:
				requestQ = append(requestQ, request)
			case worker := <-s.workerChan:
				workerQ = append(workerQ, worker)
			// 当 requestQ 或 workerQ 为空时, activeRequest 和 activeWorker 为 nil, 此 case 永远不会执行
			// 只有当真正执行了把 Request 发送给 worker chan 的操作后, 才从队列中删除数据
			case activeWorker <- activeRequest:
				requestQ = requestQ[1:]
				workerQ = workerQ[1:]
			}
		}
	}()
}
