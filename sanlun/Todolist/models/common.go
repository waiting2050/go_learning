package models

type Response struct {
	Status int `json:"status"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}

type DataList struct {
	Items interface{} `json:"items"`
	Total int64 `json:"total"`
}