package main

import (
	"github.com/dreamingtech/imoocgocrawler/frontend/controller"
	"net/http"
)

func main() {
	// curl --location 'localhost:9999/search?q=%E7%94%B7%20%E6%9C%89%E6%88%BF&from=20' \
	// --header 'Content-Type: application/json'
	// 使用空的 handler, 仅仅是为了测试
	// http.Handle("/search", controller.SearchResultHandler{})

	http.Handle("/search", controller.CreateSearchResultHandler("view/template.html"))

	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		panic(err)
	}
}
