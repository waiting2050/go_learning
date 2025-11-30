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
	
	err = dao.InitRedis()
	if err != nil {
		panic(err)
	}
	
	r := routers.SetupRouter()

	r.Run(":8080")
}