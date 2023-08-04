package controller

import (
	"github.com/dreamingtech/imoocgocrawler/engine"
	"github.com/dreamingtech/imoocgocrawler/frontend/model"
	"github.com/dreamingtech/imoocgocrawler/frontend/view"
	"github.com/olivere/elastic/v7"
	"golang.org/x/net/context"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// SearchResultHandler 把数据从 client 中取出来, 然后传给 view 进行渲染
type SearchResultHandler struct {
	view   view.SearchResultView
	client *elastic.Client
}

// CreateSearchResultHandler 初始化 SearchResultHandler
func CreateSearchResultHandler(template string) SearchResultHandler {
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
	return SearchResultHandler{
		view:   view.CreateSearchResultView(template),
		client: client,
	}
}

// ServeHTTP 处理类似请求 localhost:8888/search?q=男 已购房&from=20
func (handler SearchResultHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	q := strings.TrimSpace(request.FormValue("q"))
	// 如果在这里重写查询字符串, 重写后的查询字符串会被传递给 getSearchResult 方法
	// 然后也会被传递给前端页面, 填充到搜索框中, 每次点击都会在原有的基础上添加 `Payload.` 前缀
	// 搜索结果就不正确了, 所以应该在 getSearchResult 方法中重写查询字符串
	// q = rewriteQueryString(q)

	from, err := strconv.Atoi(request.FormValue("from"))
	if err != nil {
		from = 0
	}
	var page model.SearchResult
	page, err = handler.getSearchResult(q, from)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	err = handler.view.Render(writer, page)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
}

func (handler SearchResultHandler) getSearchResult(q string, from int) (model.SearchResult, error) {
	var result model.SearchResult

	resp, err := handler.client.
		Search("dating_profile").
		// Query(elastic.NewQueryStringQuery(q)).
		// 重写查询字符串, 在搜索前重写查询字符串, 而填充到搜索框中的查询字符串依然是原来的 q
		Query(elastic.NewQueryStringQuery(rewriteQueryString(q))).
		From(from).
		Do(context.Background())
	if err != nil {
		return result, err
	}

	result.Hits = resp.TotalHits()
	result.Start = from
	result.Query = q

	// Each 返回的是 []interface{}, 需要转换成 []engine.Item
	// 两种解决方法, 一种是把 model.SearchResult.Items 的类型改成 []interface{}
	// 再使用 Each 就不会报错了
	result.Items = resp.Each(reflect.TypeOf(engine.Item{}))

	// reflect.TypeOf(engine.Item{}) 返回的是 reflect.Type 类型
	// reflect 反射, 用于在运行时检查变量的类型和值, 例如:
	// var x float64 = 3.4
	// fmt.Println("type:", reflect.TypeOf(x))
	// fmt.Println("value:", reflect.ValueOf(x).String())
	// fmt.Println("type:", reflect.TypeOf(x).String())
	// fmt.Println("kind:", reflect.ValueOf(x).Kind())

	// 另一种是使用 type assertion
	// for _, v := range resp.Each(reflect.TypeOf(engine.Item{})) {
	// 	if item, ok := v.(engine.Item); ok {
	// 		result.Items = append(result.Items, item)
	// 	}
	// }

	// result.PrevFrom 和 result.NextFrom 要放在 result.Items 赋值之后再计算, 否则 result.Items 的长度始终为 0
	// 上一页的起始位置, 等于当前页的起始位置减去每页的条数
	result.PrevFrom = result.Start - len(result.Items)
	// 下一页的起始位置, 等于当前页的起始位置加上每页的条数
	result.NextFrom = result.Start + len(result.Items)

	// elastic.NewQueryStringQuery(q) 会把 q 中的空格替换成 +, 例如 q=男 已购房 会变成 q=男+已购房
	return result, nil
}

// rewriteQueryString 重写查询字符串
// 用户要查询年龄小于 30 岁的, 要使用 Payload.Age:(<30), 希望用户输入的是 Age:(<30)
func rewriteQueryString(q string) string {
	re := regexp.MustCompile(`([A-Z][a-z]*):`)
	return re.ReplaceAllString(q, "Payload.$1:")
}
