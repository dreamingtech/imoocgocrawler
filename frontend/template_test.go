package frontend

import (
	"github.com/dreamingtech/imoocgocrawler/engine"
	modelFront "github.com/dreamingtech/imoocgocrawler/frontend/model"
	model "github.com/dreamingtech/imoocgocrawler/model/zhenai"
	"os"
	"testing"
	"text/template"
)

func TestTemplate(t *testing.T) {
	// 因为是自己写的 template, 认为其中的语言都是合法的, 如果不合法, 直接 panic
	// templateFile := template.Must(template.ParseFiles("template.html"))
	files, err := template.ParseFiles("template.html")
	if err != nil {
		panic(err)
	}

	// 创建一个 SearchResult 对象
	page := modelFront.SearchResult{}

	// 如果不填入数据, 就会显示 template 中 else 里面的内容
	page.Hits = 123
	item := engine.Item{
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

	for i := 0; i < 10; i++ {
		page.Items = append(page.Items, item)
	}

	// 把输出在终端中打印出来
	// os.Stdout 是一个 io.Writer, 表示把输出写到标准输出, 即显示在终端
	// err = files.Execute(os.Stdout, page)

	// 把输出写到文件
	out, err := os.Create("template.test.html")
	if err != nil {
		panic(err)
	}
	err = files.Execute(out, page)
	if err != nil {
		panic(err)
	}
}
