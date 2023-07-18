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
)

func determineEncoding(reader io.Reader) encoding.Encoding {
	bytes, err := bufio.NewReader(reader).Peek(1024)
	if err != nil {
		// 如果出错, 就返回默认的 utf-8 编码
		log.Printf("error in determineEncoding: %v", err)
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e

}

func Fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// return nil, fmt.Errorf("wrong status code: %d", resp.StatusCode)
		return nil, errors.New(fmt.Sprintf("wrong status code: %d", resp.StatusCode))
	}

	e := determineEncoding(resp.Body)

	reader := transform.NewReader(resp.Body, e.NewDecoder())

	return io.ReadAll(reader)

}
