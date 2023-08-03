package model

import "github.com/dreamingtech/imoocgocrawler/engine"

type SearchResult struct {
	Hits  int
	Start int
	Items []engine.Item
}
