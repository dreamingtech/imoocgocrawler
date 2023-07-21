package scheduler

import "github.com/dreamingtech/imoocgocrawler/engine"

type SimpleScheduler struct {
	// Scheduler 为了能够保存 Request 请求, 必须要有一个 chan of Request
	workerChan chan engine.Request
}

// Submit 实现 iScheduler 中的方法
func (s *SimpleScheduler) Submit(request engine.Request) {
	// send request down to worker chan
	s.workerChan <- request
}

// ConfigureWorkerChan 把外部创建的 channel 设置为 Scheduler 的 channel
// ConfigureWorkerChan 要改变 s 中的内容, 所以要用指针类型
func (s *SimpleScheduler) ConfigureWorkerChan(c chan engine.Request) {
	s.workerChan = c
}
