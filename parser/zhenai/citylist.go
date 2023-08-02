package parser

import (
	"github.com/dreamingtech/imoocgocrawler/engine"
	"regexp"
)

const cityListRe = `<a href="([^><"]+www\.zhenai\.com/zhenghun/[0-9a-z]+)"[^>]*>([^<]+)</a>`

func ParseCityList(html []byte, _ string) engine.ParseResult {

	re := regexp.MustCompile(cityListRe)
	matches := re.FindAllSubmatch(html, -1)

	parsedResult := engine.ParseResult{}

	// 限制城市数量, 便于测试
	// limit := 10

	for _, m := range matches {

		// 把提取到的城市名字存到 Items 中
		// m[2] 的类型是 []byte,
		// 打印时, 会打印出 byte 数组的地址, 而不是字符串, 所以需要转换成字符串
		// 为方便打印时区分, 在前面加上 "City "

		// 只保存对我们有价值的 item, 后续保存到数据库中, 如用户的 Profile,
		// 就不用提取到 City 名字了, 可以使用日志来记录城市信息
		// parsedResult.Items = append(parsedResult.Items, "City "+string(m[2]))

		// 对于提取到的每一个 url, 都生成一个 Request
		parsedResult.Requests = append(parsedResult.Requests, engine.Request{
			Url: string(m[1]),
			// nil 可以编译通过, 但不能调用, 所以这里不能写成 nil, 而是定义一个空的解析函数
			// ParserFunc: nil,
			ParserFunc: ParseCity,
		})

		// limit--
		// if limit <= 0 {
		// 	break
		// }
	}
	return parsedResult

}
