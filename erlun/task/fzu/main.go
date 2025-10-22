package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// 存通知的结构，随便定义的
type Notice struct {
	Title   string
	Author  string
	Date    string
	Content string
	URL     string
}

// 全局的client，省得每次创建
var client http.Client

func main() {
	start := time.Now()

	// 先爬链接，从290到291页
	urlList := getLinks(290, 291)
	fmt.Printf("一共拿到 %d 个链接\n", len(urlList))

	// 再爬详情
	noticeList := getDetails(urlList)
	fmt.Printf("爬到了 %d 条通知\n", len(noticeList))

	// 随便打几个看看
	for i, n := range noticeList {
		if i >= 5 {
			break
		}
		fmt.Printf("\n%d. %s\n", i+1, n.Title)
		fmt.Printf("   时间: %s\n", n.Date)
		fmt.Printf("   来源: %s\n", n.Author)
		fmt.Printf("   内容: %.100s...\n", n.Content)
	}

	fmt.Printf("\n完事了，花了: %s\n", time.Since(start))
}

// 爬取指定页数的所有链接
func getLinks(startPage, endPage int) []string {
	var allUrls []string
	var lock sync.Mutex
	var wg sync.WaitGroup

	// 循环每页
	for page := startPage; page <= endPage; page++ {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			// 爬当前页的链接
			pageUrls := getOnePageLinks(p)
			// 加锁存起来
			lock.Lock()
			allUrls = append(allUrls, pageUrls...)
			lock.Unlock()
		}(page)
		// 慢点爬，别被封了
		time.Sleep(100 * time.Millisecond)
	}

	wg.Wait()
	// 去重一下，可能有重复的
	return removeSame(allUrls)
}

// 爬单页的链接
func getOnePageLinks(page int) []string {
	// 拼页面url
	url := fmt.Sprintf("https://info22-443.webvpn.fzu.edu.cn/lm_list.jsp?totalpage=1093&PAGENUM=%d&urltype=tree.TreeTempUrl&wbtreeid=1460", page)
	
	// 建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("第%d页建请求错了: %v\n", page, err)
		return nil
	}
	// 设置头，抄的浏览器的
	setHeaders(req)

	// 发请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("第%d页请求失败: %v\n", page, err)
		return nil
	}
	defer resp.Body.Close()

	// 解析html
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Printf("第%d页解析错了: %v\n", page, err)
		return nil
	}

	var urls []string
	// 找li下面的第二个a标签，试了下这个选择器能用
	doc.Find("li a:nth-child(2)").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		// 只留content.jsp的链接，别的不要
		if ok && strings.Contains(href, "content.jsp") {
			// 补全url
			fullUrl := fixUrl(href)
			urls = append(urls, fullUrl)
		}
	})

	fmt.Printf("第%d页找到了%d个链接\n", page, len(urls))
	return urls
}

// 爬详情页
func getDetails(urls []string) []Notice {
	var notices []Notice
	var lock sync.Mutex
	var wg sync.WaitGroup

	// 每个链接开个协程爬
	for _, u := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			// 爬单个详情
			notice, err := getOneDetail(url)
			if err == nil {
				lock.Lock()
				notices = append(notices, notice)
				lock.Unlock()
			} else {
				fmt.Printf("爬%s错了: %v\n", url, err)
			}
		}(u)
		// 慢点，不然容易挂
		time.Sleep(50 * time.Millisecond)
	}

	wg.Wait()
	return notices
}

// 爬单个详情页
func getOneDetail(url string) (Notice, error) {
	var n Notice
	n.URL = url

	// 建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return n, err
	}
	setHeaders(req)

	// 发请求
	resp, err := client.Do(req)
	if err != nil {
		return n, err
	}
	defer resp.Body.Close()

	// 解析
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return n, err
	}

	// 拿标题，试了下conth1这个class是对的
	n.Title = strings.TrimSpace(doc.Find("div.conth1").Text())

	// 拿日期和来源，在conthsj里
	infoHtml, _ := doc.Find("div.conthsj").Html()
	if infoHtml != "" {
		// 用正则抓日期，格式是xxxx-xx-xx
		dateMatch := regexp.MustCompile(`日期：\s*(\d{4}-\d{2}-\d{2})`).FindStringSubmatch(infoHtml)
		if len(dateMatch) > 1 {
			n.Date = dateMatch[1]
		}
		// 抓信息来源
		authorMatch := regexp.MustCompile(`信息来源：\s*([^&<]+)`).FindStringSubmatch(infoHtml)
		if len(authorMatch) > 1 {
			n.Author = authorMatch[1]
		}
	}

	// 拿正文，看了下页面里是v_news_content这个class
	n.Content = strings.TrimSpace(doc.Find(".v_news_content").Text())

	return n, nil
}

// 设置请求头，抄的浏览器的，改了点
func setHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Cookie", "_gscu_1331749010=27499092kmwl5c23; Ecp_ClientId=c241216221301948425; JSESSIONID=751C156BD637917BCB0314B85C450319; _webvpn_key=eyJhbGciOiJIUzI1NiJ9.eyJ1c2VyIjoiMTAyNDAwMzE0IiwiZ3JvdXBzIjpbNl0sImlhdCI6MTc2MTA2MDUxNywiZXhwIjoxNzYxMTQ2OTE3fQ.ifdGNM02vw4tohC2hgrKEQckx_7ypyOeaybtCHV7K-U; webvpn_username=102400314%7C1761060517%7C27093513163d98925e9064688488bb651eaad2bb")
}

// 补全url，相对路径转绝对
func fixUrl(href string) string {
	if strings.HasPrefix(href, "http") {
		return href
	}
	return "https://info22-443.webvpn.fzu.edu.cn/" + href
}

// 去重，用map简单处理下
func removeSame(urls []string) []string {
	seen := make(map[string]bool)
	var unique []string
	for _, u := range urls {
		if !seen[u] {
			seen[u] = true
			unique = append(unique, u)
		}
	}
	return unique
}