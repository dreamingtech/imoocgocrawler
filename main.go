package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
	"io"
	"net/http"
	"regexp"
)

func determineEncoding(reader io.Reader) encoding.Encoding {
	bytes, err := bufio.NewReader(reader).Peek(1024)
	if err != nil {
		panic(err)
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e

}

func getCites() {
	// const url = "http://www.zhenai.com/zhenghun"
	const url = "http://localhost:8080/mock/www.zhenai.com/zhenghun"

	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: status code", resp.StatusCode)
		return
	}

	e := determineEncoding(resp.Body)

	reader := transform.NewReader(resp.Body, e.NewDecoder())

	html, err := io.ReadAll(reader)

	if err != nil {
		panic(err)
	}
	// fmt.Printf("%s\n", html)

	printCityList(html)

}

func printCityList(html []byte) {
	// <a href="http://ww
	// w.zhenai.com/zhenghun/aba" data-v-602e7f5e>阿坝</a>
	// <a href="http://www.zhenai.com/zhenghun/aba" data-v-602e7f5e>阿坝</a>
	// [^>]* 表示匹配 0 个或多个非 > 的字符
	// [^><"] 表示匹配 0 个或多个非 > < " 的字符
	// todo 20230718 版本中, www.zhenai.com/ 中有空格, 为什么也能匹配到 ??
	// re := regexp.MustCompile(`<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)"[^>]*>([^<]+)</a>`)
	re := regexp.MustCompile(`<a href="([^><"]+www\.zhenai\.com/zhenghun/[0-9a-z]+)"[^>]*>([^<]+)</a>`)
	// [][][]byte, 最后一个 []byte 相当于一个字符串, [][]string, 第一个 [] 是所有匹配的字符串, []string 是子匹配
	matches := re.FindAllSubmatch(html, -1)
	// fmt.Printf("%s\n", matches)
	fmt.Printf("total: %d\n", len(matches))

	for _, m := range matches {
		fmt.Printf("City: %s, URL: %s\n", m[2], m[1])
	}
}

func main() {
	getCites()
}
