package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func HttpGet(url string) (res string, err error) {
	resp, err1 := http.Get(url)
	if err1 != nil {
		err = err1
		return
	}
	defer resp.Body.Close()

	// 循环读入网页数据，传出给调用者
	buf := make([]byte, 4096)
	for {
		n, err2 := resp.Body.Read(buf)
		if n == 0 {
			fmt.Println("读取网页完成")
			break
		}
		if err2 != nil && err2 != io.EOF {
			err = err2
		}
		res += string(buf[:n])
	}

	return
}

func work(st, ed int) {
	fmt.Printf("正在爬取%d到%d页....\n", st, ed)

	for i := st; i <= ed; i++ {
		url := "https://tieba.baidu.com/f?kw=%E4%B8%89%E8%A7%92%E6%B4%B2&ie=utf-8&pn=" + strconv.Itoa((i - 1) * 50)
		res, err := HttpGet(url)
		if err != nil {
			fmt.Println("HttpGet error:", err)
			continue
		}
		
		f, err := os.Create("第" + strconv.Itoa(i) + "页" + ".html")
		if err != nil {
			fmt.Println("HttpGet error:", err)
			continue
		}
		f.WriteString(res)
		f.Close()
	}
}

func main() {

	var st, ed int

	fmt.Print("请输入爬取的起始页()>=1):")
	fmt.Scan(&st)
	fmt.Print("请输入爬取的终止页()>=start):")
	fmt.Scan(&ed)

	work(st, ed)
}
