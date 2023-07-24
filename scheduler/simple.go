package scheduler

import "github.com/dreamingtech/imoocgocrawler/engine"

type SimpleScheduler struct {
	// Scheduler 为了能够保存 Request 请求, 必须要有一个 chan of Request
	workerChan chan engine.Request
}

func (s *SimpleScheduler) GetWorkerChan() chan engine.Request {
	return s.workerChan
}

func (s *SimpleScheduler) WorkerReady(requests chan engine.Request) {
}

func (s *SimpleScheduler) Run() {
	// 在 scheduler.Run 中创建一个所有 worker 公用的 channel
	s.workerChan = make(chan engine.Request)
}

// Submit 实现 iScheduler 中的方法
func (s *SimpleScheduler) Submit(request engine.Request) {
	// send request down to worker chan
	// 为了解决 worker 和 scheduler 之间循环等待的问题,
	// 可以把向 workerChan 中提交 request 的工作放在一个 goroutine 中来完成
	go func() {
		s.workerChan <- request
	}()
}

// ConfigureWorkerChan 把外部创建的 channel 设置为 Scheduler 的 channel
// ConfigureWorkerChan 要改变 s 中的内容, 所以要用指针类型
func (s *SimpleScheduler) ConfigureWorkerChan(c chan engine.Request) {
	s.workerChan = c
}
