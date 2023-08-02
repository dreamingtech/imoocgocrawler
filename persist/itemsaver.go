package persist

import (
	"errors"
	"github.com/dreamingtech/imoocgocrawler/engine"
	"github.com/olivere/elastic/v7"
	"golang.org/x/net/context"
	"log"
)

func ItemSaver(index string) (chan engine.Item, error) {
	client, err := elastic.NewClient(
		elastic.SetURL("http://127.0.0.1:9200"),
		// Must turn off sniff in docker
		elastic.SetSniff(false),
	)
	if err != nil {
		return nil, err
	}

	out := make(chan engine.Item)
	// 真正执行保存数据的操作
	go func() {
		itemCount := 0
		for {
			item := <-out
			log.Printf("ItemSaver got item #%d: %v", itemCount, item)
			itemCount++

			err := save(client, index, item)
			if err != nil {
				log.Printf("ItemSaver: error saving item %v: %v", item, err)
				continue
			}
		}
	}()
	return out, nil
}

// 可以调用 http.Post 发送请求保存数据, 也可以使用 elasticsearch 的 api 保存数据
// 为了能够测试已经保存的数据是否可以被取出来并被解析为 profile 的结构, 要返回保存到 es 时的 id
func save(client *elastic.Client, index string, item engine.Item) error {

	// client.Index().Index("dating_profile").Type("zhenai").Id("abc").BodyJson(item).Do(context.Background())
	// 不设置 id 时, es 会自动创建 id
	// todo 'Type' is deprecated
	// Index 由程序配置来指定, 而 Type 和 Id 则由 Parser 解析出来的 item.
	// 即相当于 数据库 Index 是确定的, 而数据表 Type 和 表的主键 Id 则是从响应中解析出来的

	// 必须要有 Type
	if item.Type == "" {
		return errors.New("must supply Type")
	}

	indexService := client.Index().
		Index(index).
		Type(item.Type).
		Id(item.Id).
		BodyJson(item)

	// 如果不设置 id, es 会自动创建 id
	if item.Id != "" {
		indexService.Id(item.Id)
	}

	resp, err := indexService.Do(context.Background())

	log.Printf("save item resp test %+v", resp)

	if err != nil {
		return err
	}

	return nil

}
