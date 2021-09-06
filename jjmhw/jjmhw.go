package main

import (
	zipUtils "Spider/utils"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	temp := "./out/temp"
	chapter := temp
	zipDir := "./out/zip"
	zipFile := "x.zip"
	bookId := "514"
	bookHome := "https://www.jjmhw.cc/book/" + bookId

	// Instantiate default collector
	c := colly.NewCollector()

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnRequest(func(rq *colly.Request) {
		rq.Headers.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		rq.Headers.Add("Accept-Encoding", "gzip, deflate")
		rq.Headers.Add("Accept-Language", "zh,en-US;q=0.9,en;q=0.8,zh-TW;q=0.7,zh-CN;q=0.6")
		rq.Headers.Add("Cache-Control", "no-cache")
		rq.Headers.Add("Connection", "keep-alive")
		rq.Headers.Add("Host", "www.jjmhw.cc")
		rq.Headers.Add("Pragma", "no-cache")
		rq.Headers.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36")

	})

	c.OnResponse(func(r *colly.Response) {
		//fmt.Println(string(r.Body))

		if strings.Contains(r.Headers.Get("content-type"), "image") {
			fileName := path.Base(r.Request.URL.String())
			f, err := os.Create(chapter + "/" + fileName)
			if err != nil {
				panic(err)
			}
			io.Copy(f, bytes.NewReader(r.Body))
		}

	})

	// 标题
	c.OnHTML("div[class=\"info\"] h1", func(e *colly.HTMLElement) {
		fmt.Printf("返回值: %s \n", e.Text)
		text := strings.Replace(e.Text, " ", "", -1)  // 去 空格
		text = strings.Replace(e.Text, "?", "", -1)	// 去 ？
		temp =  temp + "/" + text
		os.MkdirAll(temp, 0711)
		zipDir = zipDir + "/" + text
		os.MkdirAll(zipDir, 0711)
	})

	// 章节列表
	c.OnHTML("ul#detail-list-select > li > a", func(e *colly.HTMLElement) {
		chapterTitle := e.Text
		if chapterTitle != "" {
			println("章节==：", chapterTitle)
			url := e.Attr("href")

			chapter = temp + "/" + chapterTitle
			os.MkdirAll(chapter, 0711)

			chapterUrl := "https://www.jjmhw.cc" + url
			c.Visit(chapterUrl)
		}
	})

	// img 列表
	c.OnHTML("div[class=\"comicpage\"] > div > img", func(e *colly.HTMLElement) {
		imgUrl := e.Attr("data-original")

		c.Visit(imgUrl)
	})

	// Start scraping on https://hackerspaces.org
	c.Visit(bookHome)

	// 压缩
	items, _ := ioutil.ReadDir(temp)
	for _, item := range items {
		if item.IsDir() {
			src := temp + "/" + item.Name()
			zipFile = item.Name() + ".cbz"
			dest := zipDir + "/" + zipFile
			zipUtils.CompressFile(src, dest)
		}
	}


}