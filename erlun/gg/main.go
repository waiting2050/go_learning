package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	Client http.Client
	wg     sync.WaitGroup
)

type NoticeLink struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Page  int    `json:"page"`
}

type NoticeDetail struct {
	Title   string `json:"title"`
	Author  string `json:"author"`
	Date    string `json:"date"`
	Content string `json:"content"`
	URL     string `json:"url"`
}

func main() {
	start := time.Now()
	
	// 创建channel来收集结果
	results := make(chan NoticeLink, 1000)
	
	// 启动收集器goroutine
	var allLinks []NoticeLink
	go func() {
		for link := range results {
			allLinks = append(allLinks, link)
		}
	}()
	
	// 爬取290-410页
	for i := 290; i <= 291; i++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			SpiderPage(page, results)
		}(i)
		
		// 添加延迟，避免请求过快
		time.Sleep(100 * time.Millisecond)
	}
	
	// 等待所有爬取完成
	wg.Wait()
	close(results)
	
	fmt.Printf("共爬取到 %d 个通知链接\n", len(allLinks))
	
	// 去重处理
	uniqueLinks := removeDuplicateLinks(allLinks)
	fmt.Printf("去重后剩余 %d 个通知链接\n", len(uniqueLinks))
	
	// 现在爬取每个详情页
	var allDetails []NoticeDetail
	detailWG := sync.WaitGroup{}
	detailsChan := make(chan NoticeDetail, 100)
	
	// 启动详情收集器
	go func() {
		for detail := range detailsChan {
			allDetails = append(allDetails, detail)
		}
	}()
	
	// 并发爬取详情页
	for _, link := range uniqueLinks {
		detailWG.Add(1)
		go func(link NoticeLink) {
			defer detailWG.Done()
			detail, err := SpiderDetail(link.URL)
			if err != nil {
				log.Printf("爬取详情页失败 %s: %v", link.URL, err)
				return
			}
			detail.Title = link.Title
			detail.URL = link.URL
			detailsChan <- detail
		}(link)
		
		time.Sleep(50 * time.Millisecond) // 降低请求频率
	}
	
	detailWG.Wait()
	close(detailsChan)
	
	// 打印结果
	fmt.Printf("\n详情页爬取完成，共 %d 个通知\n", len(allDetails))
	for i, detail := range allDetails {
		if i < 5 { // 只打印前5个作为示例
			fmt.Printf("\n%d. 标题: %s\n", i+1, detail.Title)
			fmt.Printf("   发布时间: %s\n", detail.Date)
			fmt.Printf("   发布单位: %s\n", detail.Author)
			fmt.Printf("   内容预览: %.100s...\n", detail.Content)
			fmt.Printf("   链接: %s\n", detail.URL)
		}
	}
	
	elapsed := time.Since(start)
	fmt.Printf("\n全部完成，耗时: %s\n", elapsed)
}

// 去重函数 - 根据URL去重
func removeDuplicateLinks(links []NoticeLink) []NoticeLink {
	seen := make(map[string]bool)
	var unique []NoticeLink
	
	for _, link := range links {
		if !seen[link.URL] {
			seen[link.URL] = true
			unique = append(unique, link)
		}
	}
	
	return unique
}

func SpiderPage(page int, results chan<- NoticeLink) {
	url := "https://info22-443.webvpn.fzu.edu.cn/lm_list.jsp?totalpage=1093&PAGENUM=" + strconv.Itoa(page) + "&urltype=tree.TreeTempUrl&wbtreeid=1460"
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("第%d页创建请求失败: %v", page, err)
		return
	}
	
	req.Header.Set("User-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36 Edg/141.0.0.0")
	req.Header.Set("Cookie", "_gscu_1331749010=27499092kmwl5c23; Ecp_ClientId=c241216221301948425; JSESSIONID=751C156BD637917BCB0314B85C450319; _webvpn_key=eyJhbGciOiJIUzI1NiJ9.eyJ1c2VyIjoiMTAyNDAwMzE0IiwiZ3JvdXBzIjpbNl0sImlhdCI6MTc2MTA2MDUxNywiZXhwIjoxNzYxMTQ2OTE3fQ.ifdGNM02vw4tohC2hgrKEQckx_7ypyOeaybtCHV7K-U; webvpn_username=102400314%7C1761060517%7C27093513163d98925e9064688488bb651eaad2bb")

	resp, err := Client.Do(req)
	if err != nil {
		log.Printf("第%d页请求失败: %v", page, err)
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("第%d页解析失败: %v", page, err)
		return
	}

	count := 0
	// 查找所有包含通知的li
	doc.Find("li").Each(func(i int, li *goquery.Selection) {
		// 查找每个li中的第二个a标签（通知标题链接）
		link := li.Find("a:nth-child(2)")
		if link.Length() > 0 {
			text := strings.TrimSpace(link.Text())
			href, exists := link.Attr("href")
			
			// 过滤有效的通知链接
			if exists && href != "" && text != "" && strings.Contains(href, "content.jsp") {
				fullURL := buildFullURL(href)
				
				results <- NoticeLink{
					Title: text,
					URL:   fullURL,
					Page:  page,
				}
				count++
			}
		}
	})
	
	if count == 0 {
		log.Printf("第%d页未找到任何通知链接", page)
	} else {
		log.Printf("第%d页找到 %d 个通知链接", page, count)
	}
}

func SpiderDetail(url string) (NoticeDetail, error) {
	detail := NoticeDetail{}
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return detail, err
	}
	
	req.Header.Set("User-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36 Edg/141.0.0.0")
	req.Header.Set("Cookie", "_gscu_1331749010=27499092kmwl5c23; Ecp_ClientId=c241216221301948425; JSESSIONID=751C156BD637917BCB0314B85C450319; _webvpn_key=eyJhbGciOiJIUzI1NiJ9.eyJ1c2VyIjoiMTAyNDAwMzE0IiwiZ3JvdXBzIjpbNl0sImlhdCI6MTc2MTA2MDUxNywiZXhwIjoxNzYxMTQ2OTE3fQ.ifdGNM02vw4tohC2hgrKEQckx_7ypyOeaybtCHV7K-U; webvpn_username=102400314%7C1761060517%7C27093513163d98925e9064688488bb651eaad2bb")

	resp, err := Client.Do(req)
	if err != nil {
		return detail, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return detail, err
	}

	// 提取标题 - 从 conth1 类中提取
	detail.Title = strings.TrimSpace(doc.Find("div.conth1").Text())

	// 提取日期、作者 - 从 conthsj 类中提取
	infoHTML, _ := doc.Find("div.conthsj").Html()
	if infoHTML != "" {
		// 提取日期
		dateRe := regexp.MustCompile(`日期：\s*(\d{4}-\d{2}-\d{2})`)
		dateMatch := dateRe.FindStringSubmatch(infoHTML)
		if len(dateMatch) > 1 {
			detail.Date = strings.TrimSpace(dateMatch[1])
		}
		
		// 提取作者/信息来源
		authorRe := regexp.MustCompile(`信息来源：\s*([^&<]+)`)
		authorMatch := authorRe.FindStringSubmatch(infoHTML)
		if len(authorMatch) > 1 {
			detail.Author = strings.TrimSpace(authorMatch[1])
		}
	}

	// 提取正文内容
	contentSelectors := []string{".v_news_content", ".content", "#content", "div[align='left']", "div[style*='line-height']"}
	for _, selector := range contentSelectors {
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			if text != "" && len(text) > len(detail.Content) && !strings.Contains(text, "日期：") {
				detail.Content = text
			}
		})
		if detail.Content != "" {
			break
		}
	}

	return detail, nil
}

func buildFullURL(href string) string {
	if href == "" {
		return ""
	}
	// 如果已经是完整URL，直接返回
	if strings.HasPrefix(href, "http") {
		return href
	}
	// 如果是相对路径，添加前缀
	return "https://info22-443.webvpn.fzu.edu.cn/" + href
}