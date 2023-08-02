package persist

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dreamingtech/imoocgocrawler/engine"
	model "github.com/dreamingtech/imoocgocrawler/model/zhenai"
	"github.com/olivere/elastic/v7"
	"testing"
)

func TestItemSaver(t *testing.T) {

	expected := engine.Item{
		Url:  "http://album.zhenai.com/u/108906739",
		Type: "zhenai",
		Id:   "108906739",
		Payload: model.Profile{
			Age:        34,
			Height:     162,
			Weight:     57,
			Income:     "3001-5000元",
			Gender:     "女",
			Name:       "安静的雪",
			Xingzuo:    "牡羊座",
			Occupation: "人事/行政",
			Marriage:   "离异",
			House:      "已购房",
			Hokou:      "山东菏泽",
			Education:  "大学本科",
			Car:        "未购车",
		},
	}

	client, err := elastic.NewClient(
		elastic.SetURL("http://127.0.0.1:9200"),
		elastic.SetSniff(false),
	)
	if err != nil {
		panic(err)
	}

	const index = "dating_profile_test"
	// 1. save item
	err = save(client, index, expected)
	if err != nil {
		panic(err)
	}

	// 测试中依赖外界环境, 如果没有启动 es docker, 测试就会失败, 可以使用 docker go client 来启动 docker es

	// todo Try to start up elasticsearch using docker go client
	// 2. fetch item
	// 从 es 中获取数据

	resp, err := client.Get().
		Index(index).
		Type(expected.Type).Id(expected.Id).
		Do(context.Background())

	if err != nil {
		panic(err)
	}
	// resp test &{Index:dating_profile Type:zhenai Id:YjYGmokB60J7JRB5eV16
	// Uid: Routing: Parent: Version:0xc00001ae00 SeqNo:0xc00001ae08 PrimaryTerm:0xc00001ae18
	// Source:[123 34 110 ... 232 180 173 232 189 166 34 125] Found:true Fields:map[] Error:<nil>}
	// fmt.Printf("resp test %+v", resp)
	// 返回的数据在 source 中
	// resp test {"name":"安静的雪","gender":"女","age":34,"height":162,"weight":57,"income":"3001-5000元",
	// "marriage":"离异","education":"大学本科","occupation":"人事/行政","hokou":"山东菏泽","xingzuo":"牡羊座","house":"已购房","car":"
	fmt.Printf("resp test %s\n", resp.Source)

	var actual engine.Item
	err = json.Unmarshal(resp.Source, &actual)
	if err != nil {
		panic(err)
	}

	/*
		2023/08/01 09:45:00 save item resp test &{Index:dating_profile Type:zhenai Id:108906739 Version:1 Result:created Shards:0xc00006a240 SeqNo:0 PrimaryTerm:0 Status:0 ForcedRefresh:false}
		resp test {"Url":"http://album.zhenai.com/u/108906739","Type":"zhenai","Id":"108906739","Payload":{"url":"","id":"","name":"安静的雪","gender":"女","age":34,"height":162,"weight":57,"income":"3001-5000元","marriage":"离异","education":"大学本科","occupation":"人事/行政","hokou":"山东菏泽","xingzuo":"牡羊座","house":"已购房","car":"未购车"}}
		--- FAIL: TestItemSaver (0.06s)
		    itemsaver_test.go:79: got {http://album.zhenai.com/u/108906739 zhenai 108906739 map[age:34 car:未购车 education:大学本科 gender:女 height:162 hokou:山东菏泽 house:已购房 id: income:3001-5000元 marriage:离异 name:安静的雪 occupation:人事/行政 url: weight:57 xingzuo:牡羊座]}, expected {http://album.zhenai.com/u/108906739 zhenai 108906739 {  安静的雪 女 34 162 57 3001-5000元 离异 大学本科 人事/行政 山东菏泽 牡羊座 已购房 未购车}}
		FAIL
		exit status 1
		FAIL    github.com/dreamingtech/imoocgocrawler/persist  0.061s
	*/

	// 3. verify fetched item

	// 把 Payload 从 map 转换为 Profile
	actualProfile, _ := model.FromJsonObj(actual.Payload)
	actual.Payload = actualProfile

	if actual != expected {
		t.Errorf("got %v, expected %v", actual, expected)
	} else {
		fmt.Println("actual equals expected, pass")
	}
}
