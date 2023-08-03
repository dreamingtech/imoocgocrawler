package controller

import (
	"fmt"
	"github.com/dreamingtech/imoocgocrawler/frontend/view"
	"github.com/olivere/elastic/v7"
	"net/http"
	"strconv"
	"strings"
)

// SearchResultHandler 把数据从 client 中取出来, 然后传给 view 进行渲染
type SearchResultHandler struct {
	view   view.SearchResultView
	client *elastic.Client
}

// ServeHTTP 处理类似请求 localhost:8888/search?q=男 已购房&from=20
func (handler SearchResultHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	q := strings.TrimSpace(request.FormValue("q"))
	from, err := strconv.Atoi(request.FormValue("from"))
	if err != nil {
		from = 0
	}
	fmt.Fprintf(writer, "q=%s, from=%d", q, from)
}
