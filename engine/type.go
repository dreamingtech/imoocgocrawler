package engine

// ParserFunc 是一个公共的解析函数类型, url 也是一个公共的字段,
// 不只在 ProfileParser 中有用, 其它所有的 Parser 都有可能会用到
// 所以把 url 添加到公共的解析函数类型中
type ParserFunc func(contents []byte, url string) ParseResult

// Request 请求对象的封装
type Request struct {
	Url        string
	ParserFunc ParserFunc
}

// ParseResult 用于存储解析后的数据
type ParseResult struct {
	Requests []Request
	// todo 仿照 scrapy.item 对城市和用户分别定义不同的 Item
	Items []Item
}

// Item 用于存储解析后的数据
// Url, Type, Id, 是所有提取到的数据都有的字段, Payload 提取到的数据中不同的字段
type Item struct {
	Url     string
	Type    string
	Id      string
	Payload interface{}
}

// NilParser 空的解析函数
func NilParser([]byte) ParseResult {
	return ParseResult{}
}
