package main

import (
	"log"
	"silun/pkg/config"
	"silun/pkg/database"
	"silun/pkg/handler"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	cfg := config.Load()

	db, err := database.Init(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close(db)

	h := server.Default(
		server.WithHostPorts(":8888"),
	)

	handler.RegisterRoutes(h, db)

	h.Spin()
}
