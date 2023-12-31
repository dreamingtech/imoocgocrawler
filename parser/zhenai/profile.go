package parser

import (
	"github.com/dreamingtech/imoocgocrawler/engine"
	model "github.com/dreamingtech/imoocgocrawler/model/zhenai"
	"regexp"
	"strconv"
)

var ageRe = regexp.MustCompile(`<td><span class="label">年龄：</span>(\d+)岁</td>`)
var heightRe = regexp.MustCompile(`<td><span class="label">身高：</span>(\d+)CM</td>`)
var weightRe = regexp.MustCompile(`<td><span class="label">体重：</span><span field="">(\d+)KG</span></td>`)
var incomeRe = regexp.MustCompile(`<td><span class="label">月收入：</span>([^<]+)</td>`)
var genderRe = regexp.MustCompile(`<td><span class="label">性别：</span><span field="">([^<]+)</span></td>`)
var xingzuoRe = regexp.MustCompile(`<td><span class="label">星座：</span><span field="">([^<]+)</span></td>`)
var marriageRe = regexp.MustCompile(`<td><span class="label">婚况：</span>([^<]+)</td>`)
var educationRe = regexp.MustCompile(`<td><span class="label">学历：</span>([^<]+)</td>`)
var occupationRe = regexp.MustCompile(`<td><span class="label">职业：</span><span field="">([^<]+)</span></td>`)
var hokouRe = regexp.MustCompile(`<td><span class="label">籍贯：</span>([^<]+)</td>`)
var houseRe = regexp.MustCompile(`<td><span class="label">住房条件：</span><span field="">([^<]+)</span></td>`)
var carRe = regexp.MustCompile(`<td><span class="label">是否购车：</span><span field="">([^<]+)</span></td>`)

// 猜你喜欢
var guessRe = regexp.MustCompile(`<a class="exp-user-name"[^>]*href="([^<>"]+album.zhenai.com/u/[\d]+)"[^>]*>([^<]+)</a>`)
var idUrlRe = regexp.MustCompile(`[^><]*album.zhenai.com/u/([\d]+)`)

func extractString(contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents)

	if len(match) >= 2 {
		return string(match[1])
	} else {
		return ""
	}
}

// ParseProfile 解析用户信息
func ParseProfile(html []byte, url, name string) engine.ParseResult {

	profile := model.Profile{}
	profile.Name = name

	if age, err := strconv.Atoi(extractString(html, ageRe)); err == nil {
		profile.Age = age
	}

	if height, err := strconv.Atoi(extractString(html, heightRe)); err == nil {
		profile.Height = height
	}

	if weight, err := strconv.Atoi(extractString(html, weightRe)); err == nil {
		profile.Weight = weight
	}

	profile.Income = extractString(html, incomeRe)
	profile.Gender = extractString(html, genderRe)
	profile.Car = extractString(html, carRe)
	profile.Education = extractString(html, educationRe)
	profile.Hokou = extractString(html, hokouRe)
	profile.House = extractString(html, houseRe)
	profile.Marriage = extractString(html, marriageRe)
	profile.Occupation = extractString(html, occupationRe)
	profile.Xingzuo = extractString(html, xingzuoRe)

	matches := guessRe.FindAllSubmatch(html, -1)

	result := engine.ParseResult{
		Items: []engine.Item{
			{
				Url:     url,
				Type:    "zhenai",
				Id:      extractString([]byte(url), idUrlRe),
				Payload: profile,
			},
		},
	}

	// 猜你喜欢提取到的用户信息添加到请求队列中
	for _, m := range matches {
		// name 的作用是防止循环中的变量被覆盖, 但现在 name 是给 ParseProfile 函数用的,
		// 函数都是值传递, 所以这里不需要再定义 name 变量了
		// name := string(m[2])
		result.Requests = append(result.Requests, engine.Request{
			Url:        string(m[1]),
			ParserFunc: ProfileParser(string(m[2])),
		})
	}

	return result
}

// ProfileParser 包装 ParseProfile, 传递 name, url 参数,
// 生成一个新的解析函数, 能够解析一段文本 []byte, 即 html 源码
func ProfileParser(name string) engine.ParserFunc {
	return func(c []byte, url string) engine.ParseResult {
		return ParseProfile(c, url, name)
	}
}
