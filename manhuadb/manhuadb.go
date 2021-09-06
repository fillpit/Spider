package main

import (
	zipUtils "Spider/utils"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	temp := "./out/temp"
	chapter := temp
	zipDir := "./out/zip"
	zipFile := "x.zip"
	bookId := "25797"
	bookHome := "https://www.manhuadb.com/manhua/" + bookId

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
		rq.Headers.Add("Host", "www.manhuadb.com")
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
	c.OnHTML("h1[class=\"comic-title\"]", func(e *colly.HTMLElement) {
		fmt.Printf("返回值: %s \n", e.Text)
		text := strings.Replace(e.Text, " ", "", -1)  // 去 空格
		text = strings.Replace(e.Text, "?", "", -1)	// 去 ？
		temp =  temp + "/" + text
		os.MkdirAll(temp, 0711)
		zipDir = zipDir + "/" + text
		os.MkdirAll(zipDir, 0711)
	})

	// 章节列表
	c.OnHTML("li[class=\"sort_div fixed-wd-num\"] > a", func(e *colly.HTMLElement) {
		chapterTitle := e.Attr("title")
		if chapterTitle != "" {
			println("章节==：", chapterTitle)
			url := e.Attr("href")

			chapter = temp + "/" + chapterTitle
			os.MkdirAll(chapter, 0711)

			chapterUrl := "https://www.manhuadb.com" + url
			c.Visit(chapterUrl)
		}
	})

	// img 列表
	c.OnHTML("div[class=\"text-center pjax-container\"] > img", func(e *colly.HTMLElement) {
		imgUrl := e.Attr("src")
		if imgUrl != "" {
			c.Visit(imgUrl)
		}
	})

	// 判断是否还有下一页
	c.OnHTML("div[class=\"container-fluid comic-detail p-0\"]", func(e *colly.HTMLElement) {
		pageHref := e.ChildAttr("li[class=\"breadcrumb-item active\"] > a", "href")
		currentPage := e.ChildText("li[class=\"breadcrumb-item active\"] > span")
		reg := regexp.MustCompile("共 ([0-9]*?) 页")

		sumPageSize := reg.FindStringSubmatch(e.ChildText("li[class=\"breadcrumb-item active\"]"))

		hrefs := strings.Split(pageHref, ".")
		sumPageSizeNum, err:=strconv.Atoi(sumPageSize[1])
		currentPageNum, err:=strconv.Atoi(currentPage)
		if err != nil {
			println(err)
		}
		if currentPageNum < sumPageSizeNum {
			p :=strconv.Itoa(currentPageNum + 1)
			next := "https://www.manhuadb.com" + hrefs[0] + "_p"+ p +"." + hrefs[1]

			c.Visit(next)
		}

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