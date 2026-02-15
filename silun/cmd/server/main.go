package main

import (
	"log"
	"silun/pkg/cache"
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

	if err := cache.InitRedis(cfg); err != nil {
		log.Printf("Warning: Failed to initialize Redis: %v", err)
	}
	defer cache.CloseRedis()

	h := server.Default(
		server.WithHostPorts("0.0.0.0:8888"),
	)

	handler.RegisterRoutes(h, db, cfg)

	h.Spin()
}
