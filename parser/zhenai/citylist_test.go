package parser

import (
	"os"
	"testing"
)

func TestParseCityList(t *testing.T) {
	// 测试时, 要尽量避免因为网络等外部因素引起的可能的错误
	// 可以把 html 保存到本地, 使用文件读取的方式, 而是不发送网络请求
	// html, err := fetcher.Fetch("https://www.zhenai.com/zhenghun")
	// html, err := fetcher.Fetch("http://localhost:8080/mock/www.zhenai.com/zhenghun")

	html, err := os.ReadFile("citylist_test_data.html")

	if err != nil {
		panic(err)
	}

	result := ParseCityList(html, "")

	// verify result
	const resultSize = 470

	if len(result.Requests) != resultSize {
		t.Errorf("result should have %d requests; but had %d", resultSize, len(result.Requests))
	}

	// ParseCityList 中已经不再提取出 item 了, 所以也不用再测试 item 了
	// if len(result.Items) != resultSize {
	// 	t.Errorf("result should have %d items; but had %d", resultSize, len(result.Items))
	// }

	expectedUrls := []string{
		"http://www.zhenai.com/zhenghun/aba",
		"http://www.zhenai.com/zhenghun/akesu",
		"http://www.zhenai.com/zhenghun/alashanmeng",
	}
	// expectedCities := []string{
	// 	"City 阿坝", "City 阿克苏", "City 阿拉善盟",
	// }

	for i, url := range expectedUrls {
		if result.Requests[i].Url != url {
			t.Errorf("expected url #%d: %s; but was %s", i, url, result.Requests[i].Url)
		}
	}

	// for i, city := range expectedCities {
	// 	// items 中的元素类型是 interface{}, 需要使用类型断言
	// 	if result.Items[i].(string) != city {
	// 		t.Errorf("expected city #%d: %s; but was %s", i, city, result.Items[i].(string))
	// 	}
	// }
}
