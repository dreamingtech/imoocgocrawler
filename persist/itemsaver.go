package persist

import (
	"fmt"
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
func save(item interface{}) {
	client, err := elastic.NewClient(
		elastic.SetURL("http://127.0.0.1:9200"),
		// Must turn off sniff in docker
		elastic.SetSniff(false),
	)
	if err != nil {
		panic(err)
	}
	// client.Index().Index("dating_profile").Type("zhenai").Id("abc").BodyJson(item).Do(context.Background())
	// 不设置 id 时, es 会自动创建 id

	resp, err := client.Index().Index("dating_profile").Type("zhenai").BodyJson(item).Do(context.Background())

	if err != nil {
		panic(err)
	}

	fmt.Printf("resp: %#v", resp)

}
