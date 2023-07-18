package main

import (
	"github.com/dreamingtech/imoocgocrawler/engine"
	parser "github.com/dreamingtech/imoocgocrawler/parser/zhenai"
)

func main() {
	engine.Run(engine.Request{
		// Url: "http://www.zhenai.com/zhenghun",
		Url:        "http://localhost:8080/mock/www.zhenai.com/zhenghun",
		ParserFunc: parser.ParseCityList,
	})
}
