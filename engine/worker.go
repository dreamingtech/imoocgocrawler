package engine

import (
	"github.com/dreamingtech/imoocgocrawler/fetcher"
	"log"
)

func doWork(r Request) (ParseResult, error) {

	log.Printf("Fetching url: %s", r.Url)

	body, err := fetcher.Fetch(r.Url)

	if err != nil {
		log.Printf("Fetcher url error. url: %s, error: %v", r.Url, err)
		return ParseResult{}, err
	}

	// 调用解析函数, 得到解析结果
	return r.ParserFunc(body, r.Url), nil
}
