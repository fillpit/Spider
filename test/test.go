package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	temp := "./out/temp"
	bookHome := "https://audio.yousxs.com:9002/mydate/%E6%81%90%E6%80%96%E7%8E%84%E5%B9%BB/%E6%97%A0%E4%BA%BA%E7%94%9F%E8%BF%98_%E6%9D%8E%E9%87%8E%E5%A2%A8/001_A.mp3?skey=17ab4e829a5a46bd915e63d9ff696328"

	// Instantiate default collector
	c := colly.NewCollector()

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnRequest(func(rq *colly.Request) {
		rq.Headers.Add("Accept", "*/*")
		rq.Headers.Add("Accept-Encoding", "identity;q=1, *;q=0")
		rq.Headers.Add("Accept-Language", "zh,en-US;q=0.9,en;q=0.8,zh-TW;q=0.7,zh-CN;q=0.6")
		rq.Headers.Add("Referer", "https://www.yousxs.com/")
		rq.Headers.Add("Sec-Fetch-Dest", "audio")
		rq.Headers.Add("Host", "audio.yousxs.com:9002")
		rq.Headers.Add("sec-ch-ua", "\"Chromium\";v=\"92\", \" Not A;Brand\";v=\"99\", \"Google Chrome\";v=\"92\"")
		rq.Headers.Add("sec-ch-ua-mobile", "?0")
		rq.Headers.Add("Sec-Fetch-Mode", "no-cors")
		rq.Headers.Add("same-site", "same-site")
		rq.Headers.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36")

	})

	c.OnResponse(func(r *colly.Response) {
		//fmt.Println(string(r.Body))

		if strings.Contains(r.Headers.Get("Content-Disposition"), "filename") {
			fileName := path.Base(strings.Split(r.Request.URL.String(), "?")[0])
			println("文件名=", fileName)
			f, err := os.Create(temp + "/" + fileName)
			if err != nil {
				panic(err)
			}
			io.Copy(f, bytes.NewReader(r.Body))
		}

	})

	// Start scraping on https://hackerspaces.org
	c.Visit(bookHome)


}