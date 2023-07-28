package persist

import (
	model "github.com/dreamingtech/imoocgocrawler/model/zhenai"
	"testing"
)

func TestItemSaver(t *testing.T) {
	profile := model.Profile{
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
	save(profile)
}
