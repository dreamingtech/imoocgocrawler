package parser

import (
	"github.com/dreamingtech/imoocgocrawler/engine"
	"regexp"
)

const cityListRe = `<a href="([^><"]+www\.zhenai\.com/zhenghun/[0-9a-z]+)"[^>]*>([^<]+)</a>`

func ParseCityList(html []byte) engine.ParseResult {

	re := regexp.MustCompile(cityListRe)
	matches := re.FindAllSubmatch(html, -1)

	parsedResult := engine.ParseResult{}

	for _, m := range matches {

		// 把提取到的城市名字存到 Items 中
		// m[2] 的类型是 []byte,
		// 打印时, 会打印出 byte 数组的地址, 而不是字符串, 所以需要转换成字符串
		parsedResult.Items = append(parsedResult.Items, string(m[2]))

		// 对于提取到的每一个 url, 都生成一个 Request
		parsedResult.Requests = append(parsedResult.Requests, engine.Request{
			Url: string(m[1]),
			// nil 可以编译通过, 但不能调用, 所以这里不能写成 nil, 而是定义一个空的解析函数
			// ParserFunc: nil,
			ParserFunc: engine.NilParser,
		})
	}
	return parsedResult

}
