package main

import "net/http"

// https://api.bilibili.com/x/v2/reply/wbi/main?oid=420981979&type=1&mode=3&pagination_str=%7B%22offset%22:%22CAESEDE4MDMyMjYyMzYyOTQwODQiAggB%22%7D&plat=1&web_location=1315875&w_rid=7c291ed09d03ba7435a435a940228f0a&wts=1760963341

func main() {
	client := http.Client{}
	req, err := http.NewRequest()
}