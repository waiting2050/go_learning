package main

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	Spider()
}

func Spider() {
	// 1．发送请求
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://movie.douban.com/top250", nil)
	if err != nil {
		fmt.Println("req err", err)
	}

	//加请求头伪造浏览器访问
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36 Edg/141.0.0.0")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("referer", "https://cn.bing.com/")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败", err)
	}

	defer resp.Body.Close()
	//2.解析网页
	docDetail, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("解析失败", err)
	}

	//3.获取节点信息
	// #content > div > div.article > ol > li:nth-child(1) > div > div.info > div.hd > a > span:nth-child(1)
	// #content > div > div.article > ol > li:nth-child(1)
	// #content > div > div.article > ol > li:nth-child(2)
	// #content > div > div.article > ol > li:nth-child(1) > div > div.pic > a > img
	// #content > div > div.article > ol > li:nth-child(1) > div > div.info > div.bd > p:nth-child(1)
	// #content > div > div.article > ol > li:nth-child(1) > div > div.info > div.bd > div > span.rating_num
	// #content > div > div.article > ol > li:nth-child(1) > div > div.info > div.bd > p.quote > span
	docDetail.Find("#content > div > div.article > ol > li").
	Each(func(i int, s *goquery.Selection) {
		title := s.Find("div > div.info > div.hd > a > span:nth-child(1)").Text()
		img := s.Find("div > div.pic > a > img")
		imgTmp, ok := img.Attr("src")
		info := s.Find("div > div.info > div.bd > p:nth-child(1)").Text()
		score := s.Find("div > div.info > div.bd > div > span.rating_num").Text()
		quote := s.Find("div > div.info > div.bd > p.quote > span").Text()

		if ok {
			fmt.Println("title", title)
			fmt.Println("imgTmp", imgTmp)
			fmt.Println("info", info)
			fmt.Println("score", score)
			fmt.Println("quote", quote)
		}
	})

	// 4.保存信息
}
