package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	temp := "./out/temp"
	bookId := "9175"
	bookHome := "https://www.yousxs.com/"+bookId+".html"

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("www.yousxs.com", "audio.yousxs.com:9002"),
	)

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
		//rq.Headers.Add("Host", "yousxs.com")
		//rq.Headers.Add("Pragma", "no-cache")
		//rq.Headers.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36")

	})

	c.OnResponse(func(r *colly.Response) {
		//fmt.Println(string(r.Body))

		if strings.Contains(r.Request.URL.String(), "player") {
			//println("==", string(r.Body))

			skey := regexp.MustCompile("skey = '(.*?)'").FindStringSubmatch(string(r.Body))
			audioUrl := regexp.MustCompile("url: '(.*?)'").FindStringSubmatch(string(r.Body))
			if len(skey) > 0 {
				url := audioUrl[1] + skey[1]

				c.Visit(url)
			}
		}

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

	// 标题
	c.OnHTML("div[class=\"col-md-7 col-xs-12 col-sm-7\"] > h3", func(e *colly.HTMLElement) {
		fmt.Printf("返回值: %s \n", e.Text)
		text := strings.Replace(e.Text, " ", "", -1)  // 去 空格
		text = strings.Replace(e.Text, "?", "", -1)	// 去 ？
		if !strings.Contains(temp, text) {
			temp =  temp + "/" + text
			os.MkdirAll(temp, 0711)
		}

	})

	// 章节列表
	c.OnHTML("div[class=\"col-md-2 col-xs-3 col-sm-2\"] > a", func(e *colly.HTMLElement) {
		chapterTitle := e.Text
		if chapterTitle != "" {
			println("章节==：", chapterTitle)
			url := e.Attr("href")

			chapterUrl := "https://www.yousxs.com" + "/" + url
			c.Visit(chapterUrl)
		}
	})

	// Start scraping on https://hackerspaces.org
	c.Visit(bookHome)


}