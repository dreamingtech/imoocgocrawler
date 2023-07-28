package persist

import (
	"github.com/olivere/elastic/v7"
	"golang.org/x/net/context"
	"log"
)

func ItemSaver() chan interface{} {
	out := make(chan interface{})
	// 真正执行保存数据的操作
	go func() {
		itemCount := 0
		for {
			item := <-out
			log.Printf("ItemSaver got item #%d: %v", itemCount, item)
			itemCount++

			save(item)
		}

	}()
	return out
}

// 可以调用 http.Post 发送请求保存数据, 也可以使用 elasticsearch 的 api 保存数据
// 为了能够测试已经保存的数据是否可以被取出来并被解析为 profile 的结构, 要返回保存到 es 时的 id
func save(item interface{}) (id string, err error) {
	client, err := elastic.NewClient(
		elastic.SetURL("http://127.0.0.1:9200"),
		// Must turn off sniff in docker
		elastic.SetSniff(false),
	)
	if err != nil {
		return "", err
	}
	// client.Index().Index("dating_profile").Type("zhenai").Id("abc").BodyJson(item).Do(context.Background())
	// 不设置 id 时, es 会自动创建 id
	// todo 'Type' is deprecated
	resp, err := client.Index().Index("dating_profile").Type("zhenai").BodyJson(item).Do(context.Background())

	if err != nil {
		return "", err
	}

	return resp.Id, nil

}
