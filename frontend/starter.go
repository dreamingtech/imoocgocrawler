package main

import (
	"github.com/dreamingtech/imoocgocrawler/frontend/controller"
	"net/http"
)

func main() {

	// 解决 js 和 css 文件无法加载的问题
	// /search 接口只处理 search 请求, 其他的请求都未处理, 所以无法加载 js 和 css 文件
	// 因为我们的 js 和 css starter 同级的 static 目录中, 所以设置目录为当前目录即可
	http.Handle("/", http.FileServer(http.Dir("view/")))

	// curl --location 'localhost:9999/search?q=%E7%94%B7%20%E6%9C%89%E6%88%BF&from=20' \
	// --header 'Content-Type: application/json'
	// 使用空的 handler, 仅仅是为了测试
	// http.Handle("/search", controller.SearchResultHandler{})

	http.Handle("/search", controller.CreateSearchResultHandler("view/template.html"))

	err := http.ListenAndServe(":3306", nil)
	if err != nil {
		panic(err)
	}
}
