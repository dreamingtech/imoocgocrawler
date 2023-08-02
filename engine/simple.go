package engine

import (
	"log"
)

type SimpleEngine struct{}

func (engine SimpleEngine) Run(seeds ...Request) {
	var requests []Request
	for _, r := range seeds {
		requests = append(requests, r)
	}

	profileCount := 0

	for len(requests) > 0 {
		// 取出来第一个 Request
		r := requests[0]
		requests = requests[1:]

		parseResult, err := doWork(r)
		if err != nil {
			continue
		}

		// 将解析结果中的 Requests 添加到 requests 列表中
		requests = append(requests, parseResult.Requests...)

		// 打印解析结果中的 Items
		for _, item := range parseResult.Items {
			log.Printf("Got profile: #%d: %v", profileCount, item)
			profileCount++
		}
	}
}
