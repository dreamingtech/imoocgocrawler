package engine

import (
	"github.com/dreamingtech/imoocgocrawler/fetcher"
	"log"
)

func Run(seeds ...Request) {
	var requests []Request
	for _, r := range seeds {
		requests = append(requests, r)
	}

	for len(requests) > 0 {
		// 取出来第一个 Request
		r := requests[0]
		requests = requests[1:]

		log.Printf("Fetching url: %s", r.Url)

		body, err := fetcher.Fetch(r.Url)

		if err != nil {
			log.Printf("Fetcher url error. url: %s, error: %v", r.Url, err)
			continue
		}

		// 调用解析函数, 得到解析结果
		parseResult := r.ParserFunc(body)

		// 将解析结果中的 Requests 添加到 requests 列表中
		requests = append(requests, parseResult.Requests...)

		// 打印解析结果中的 Items
		for _, item := range parseResult.Items {
			log.Printf("Got item: %+v", item)
		}
	}
}
