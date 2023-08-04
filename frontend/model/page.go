package model

type SearchResult struct {
	Hits  int64
	Start int
	Query string // 把查询条件也返回给前端, 填充到搜索框中
	// Items []engine.Item
	Items []interface{}
}
