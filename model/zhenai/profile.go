package model

// Profile 用户信息
type Profile struct {
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
