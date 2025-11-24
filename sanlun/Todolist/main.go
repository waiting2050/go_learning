package main

import (
	"Todolist/dao"
	"Todolist/routers"
)

func main() {
	err := dao.InitMySQL()
	if err != nil {
		panic(err)
	}

	r := routers.SerupRouter()

	r.Run(":8080")
}