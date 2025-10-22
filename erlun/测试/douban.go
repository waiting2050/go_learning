package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type MovieData struct {
	Title string `json:"title"`
	Director string `json:"director"`
	Picture string `json:"picture"`
	Actor string `json:"actor"`
	Year string `json:"year"`
	Score string `json:"score"`
	Quote string `json:"quote"`
}

func main() {
	for i := 0; i < 10; i++ {
		fmt.Printf("正在爬取第 %d 页\n", i + 1)
		Spider(strconv.Itoa(i * 25))
	}
}

func Spider(page string) {
	// 1．发送请求
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://movie.douban.com/top250?start=" + page, nil)
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
			var data MovieData
			title := s.Find("div > div.info > div.hd > a > span:nth-child(1)").Text()
			img := s.Find("div > div.pic > a > img")
			imgTmp, ok := img.Attr("src")
			info := s.Find("div > div.info > div.bd > p:nth-child(1)").Text()
			score := s.Find("div > div.info > div.bd > div > span.rating_num").Text()
			quote := s.Find("div > div.info > div.bd > p.quote > span").Text()

			if ok {
				director, actor, year := InfoSpite(info)
				data.Title = title
				data.Director = director
				data.Picture = imgTmp
				data.Actor = actor
				data.Year = year
				data.Score = score
				data.Quote = quote
				fmt.Println("data", data)
			}
		})
}

// 4.保存信息
func InfoSpite(info string) (director, actor, year string) {

	directorRe, _ := regexp.Compile(`导演:(.*)主演:`)
	director = string(directorRe.Find([]byte(info)))

	actorRe, _ := regexp.Compile(`主演:(.*)`)
	actor = string(actorRe.Find([]byte(info)))

	yearRe, _ := regexp.Compile(`(\d+)`)
	year = string(yearRe.Find([]byte(info)))

	return
}
