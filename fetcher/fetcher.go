package fetcher

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"log"
	"net/http"
	"time"
)

func determineEncoding(reader *bufio.Reader) encoding.Encoding {
	// Peek 是对 reader 的一个预读, 读取前 1024 个字节
	// 在后面 `reader := transform.NewReader(resp.Body, e.NewDecoder())` 中已经读取的 1024 个字节就不会再读取了
	// 为了解决这个问题, 可以把 bufio.NewReader 的整体作为参数传递进来
	bytes, err := reader.Peek(1024)
	if err != nil {
		// 如果出错, 就返回默认的 utf-8 编码
		log.Printf("error in determineEncoding: %v", err)
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e

}

// 添加限速器, 以免抓取过快导致被反爬
// 100 毫秒, 即 10个请求/s, 10 ms, 即 100个请求/s
var rateLimiter = time.Tick(1 * time.Millisecond)

func Fetch(url string) ([]byte, error) {
	// 多个 worker 都会调用这同一个 Fetch 完成发送请求的工作,
	// 也就会抢占式的去执行从 rateLimiter channel 中获取数据的操作
	<-rateLimiter
	// resp, err := http.Get(url)

	// 使用 http.Get 获取到的响应和浏览器中看到的响应不同, 所以要自定义请求头
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add(
		"User-Agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	resp, err := http.DefaultClient.Do(request)

	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("error in defer Body.Close(): %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		// return nil, fmt.Errorf("wrong status code: %d", resp.StatusCode)
		return nil, errors.New(fmt.Sprintf("wrong status code: %d", resp.StatusCode))
	}

	// 把 bufio.NewReader(resp.Body) 作为参数传递进去, 利用 go 语言值传递的特性,
	// 解决 peek 读取了 1024 个字节后, 后面的 transform.NewReader(resp.Body, e.NewDecoder()) 就不会再读取的问题
	bodyReader := bufio.NewReader(resp.Body)

	e := determineEncoding(bodyReader)

	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())

	return io.ReadAll(utf8Reader)

}
