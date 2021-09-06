package main

import (
	zipUtils "Spider/utils"
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
	zipDir := "./out/zip"
	zipFile := "x.zip"
	book_id := "126845"
	book_home := "https://wnacg.org/photos-index-aid-"+book_id+".html"
	book_red := "https://wnacg.org/photos-gallery-aid-"+book_id+".html"

	// Instantiate default collector
	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) {
		// 漫画信息页
		if strings.Contains(book_red, r.Request.URL.String()) {
			redBody := string(r.Body)
			reg := regexp.MustCompile("//\\w+\\.wnacg\\.org/data/\\d{4}/\\d{2}/[0-9_A-Za-z]*?\\.(jpg|png|jpeg)")

			imgList := reg.FindAllString(redBody, -1);
			fmt.Println("%q\n", imgList)
			for _, item := range imgList {
				url := "https:" + item
				c.Visit(url)
			}
		}

		if strings.Contains(r.Headers.Get("content-type"), "image") {
			fileName := path.Base(r.Request.URL.String())
			f, err := os.Create(temp + "/" + fileName)
			if err != nil {
				panic(err)
			}
			io.Copy(f, bytes.NewReader(r.Body))
		}

	})

	c.OnHTML("#bodywrap h2", func(e *colly.HTMLElement) {
		fmt.Printf("返回值: %s \n", e.Text)
		temp =  temp + "/" + e.Text
		os.MkdirAll(temp, 0711)
		zipDir = zipDir + "/" + e.Text
		os.MkdirAll(zipDir, 0711)

		zipFile = e.Text + ".cbz"
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	c.Visit(book_home)
	c.Visit(book_red)

	// 压缩
	dest := zipDir + "/" + zipFile
	zipUtils.CompressFile(temp, dest)

}