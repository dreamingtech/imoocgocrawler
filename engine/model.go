package engine

// Request 请求对象的封装
type Request struct {
	Url        string
	ParserFunc func([]byte) ParseResult
}

// ParseResult 用于存储解析后的数据
type ParseResult struct {
	Requests []Request
	// todo 仿照 scrapy.item 对城市和用户分别定义不同的 Item
	Items []interface{}
}

// NilParser 空的解析函数
func NilParser([]byte) ParseResult {
	return ParseResult{}
}
