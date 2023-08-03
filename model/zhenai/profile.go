package model

import "encoding/json"

// Profile 用户信息
// 定义了 `json:"name"` 这样的 tag, 使得在 json 序列化时, 可以将字段名转换为小写
// 但是在保存到 ElasticSearch 时, 会将字段名转换为小写, 从 ES 中取出数据填充到 template 时, 会找不到对应的字段
// 即如果定义了 `json:"name"` 这样的 tag, template 中的字段名也要小写
type JsonProfile struct {
	Name       string `json:"name"`
	Gender     string `json:"gender"`
	Age        int    `json:"age"`
	Height     int    `json:"height"`
	Weight     int    `json:"weight"`
	Income     string `json:"income"`
	Marriage   string `json:"marriage"`
	Education  string `json:"education"`
	Occupation string `json:"occupation"`
	Hokou      string `json:"hokou"`   // 籍贯, 户口
	Xingzuo    string `json:"xingzuo"` // 星座
	House      string `json:"house"`   // 是否购房
	Car        string `json:"car"`     // 是否购车
}

type Profile struct {
	Name       string
	Gender     string
	Age        int
	Height     int
	Weight     int
	Income     string
	Marriage   string
	Education  string
	Occupation string
	Hokou      string // 籍贯, 户口
	Xingzuo    string // 星座
	House      string // 是否购房
	Car        string // 是否购车
}

// FromJsonObj 从 json 对象中解析出 Profile
// 因为 engine.Item 对 Profile 进行了一层封装, 其中的 Payload 是 interface{} 类型, 所以需要转换
func FromJsonObj(o interface{}) (Profile, error) {
	var profile Profile
	s, err := json.Marshal(o)
	if err != nil {
		return profile, err
	}
	err = json.Unmarshal(s, &profile)
	return profile, err
}
