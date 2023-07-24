package parser

import (
	"github.com/dreamingtech/imoocgocrawler/engine"
	"regexp"
)

// 提取用户 url
var profileRe = regexp.MustCompile(`<a href="([^<>"]+album\.zhenai\.com/u/[0-9]+)"[^>]*>([^<]+)</a>`)

// 提取下一页 url
var cityUrlRe = regexp.MustCompile(`href="([^<>"]+www\.zhenai\.com/zhenghun/[^"]+)"`)

func ParseCity(html []byte) engine.ParseResult {

	matches := profileRe.FindAllSubmatch(html, -1)

	parsedResult := engine.ParseResult{}

	for _, m := range matches {

		// 解决每页所有用户提取到的都是最后一个用户的名字的问题
		// 循环中共用一个变量, 会导致变量的值被覆盖
		// ParseProfile 中传入的 m[2] 参数不是立即执行的, 而是等循环结束, 真正抓取到用户详情页时才执行
		// 循环结束时, m[2] 的值是最后一个用户的名字, 所以在这里需要把 m[2] 的值赋给一个变量, 然后传给 ParseProfile
		name := string(m[2])

		// 把提取到的城市名字存到 Items 中
		// m[2] 的类型是 []byte,
		// 打印时, 会打印出 byte 数组的地址, 而不是字符串, 所以需要转换成字符串
		// 为方便打印时区分, 在前面加上 "User "
		parsedResult.Items = append(parsedResult.Items, "User "+name)

		// 对于提取到的每一个 url, 都生成一个 Request
		parsedResult.Requests = append(parsedResult.Requests, engine.Request{
			Url: string(m[1]),
			// nil 可以编译通过, 但不能调用, 所以这里不能写成 nil, 而是定义一个空的解析函数
			// ParserFunc: nil,
			// 把列表页提取到的 Name 传递给下一层解析函数 ParseProfile, 而不是在详情页中解析用户的 Name
			// 但 Request.ParserFunc 中没有 Name 参数, 所以这里使用函数式编程,
			// 用一个匿名函数把 ParseProfile 进行一层包装, 并把 Name 传递给 ParseProfile
			ParserFunc: func(bytes []byte) engine.ParseResult {
				// 会出现所有用户的名字都是最后一个用户的名字的问题
				// return ParseProfile(bytes, string(m[2]))
				// return ParseProfile(bytes, name)
				// todo 测试, 不抓取用户详情页
				return engine.NilParser(bytes)
			},
		})
	}

	// 提取下一页 url, 页脚中的类似页面
	matches = cityUrlRe.FindAllSubmatch(html, -1)
	for _, m := range matches {
		parsedResult.Requests = append(parsedResult.Requests, engine.Request{
			Url:        string(m[1]),
			ParserFunc: ParseCity,
		})
	}

	return parsedResult

}
