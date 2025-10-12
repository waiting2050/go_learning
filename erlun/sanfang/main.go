// 步骤：
// 1、go mod init 项目名称
// 2、配置第三方包
// 3、go mod tidy 下载依赖

package main

import (
	"fmt"

	"github.com/shopspring/decimal"
)

func main() {

	price, err := decimal.NewFromString("136.02")
	if err != nil {
		panic(err)
	}
	fmt.Println(price)
}