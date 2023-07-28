package persist

import (
	"context"
	"encoding/json"
	"fmt"
	model "github.com/dreamingtech/imoocgocrawler/model/zhenai"
	"github.com/olivere/elastic/v7"
	"testing"
)

func TestItemSaver(t *testing.T) {
	expected := model.Profile{
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
	}
	id, err := save(expected)
	if err != nil {
		panic(err)
	}

	// 测试中依赖外界环境, 如果没有启动 es docker, 测试就会失败, 可以使用 docker go client 来启动 docker es
	// todo Try to start up elasticsearch using docker go client
	// 从 es 中获取数据
	client, err := elastic.NewClient(
		elastic.SetURL("http://127.0.0.1:9200"),
		elastic.SetSniff(false),
	)
	if err != nil {
		panic(err)
	}
	resp, err := client.Get().Index("dating_profile").Type("zhenai").Id(id).Do(context.Background())

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

	var actual model.Profile
	err = json.Unmarshal(resp.Source, &actual)
	if err != nil {
		panic(err)
	}

	if actual != expected {
		t.Errorf("got %v, expected %v", actual, expected)
	} else {
		fmt.Println("actual equals expected, pass")
	}
}
