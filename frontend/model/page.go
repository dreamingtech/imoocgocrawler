package model

type SearchResult struct {
	Hits     int64
	Start    int
	Query    string // 把查询条件也返回给前端, 填充到搜索框中
	PrevFrom int    // 前一页的起始位置
	NextFrom int    // 下一页的起始位置
	// Items []engine.Item
	Items []interface{}
}
