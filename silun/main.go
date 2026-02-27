package main

import (
	"log"
	"silun/biz/cache"
	"silun/biz/model"
	"silun/biz/auth"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	cfg := model.LoadConfig()

	db, err := model.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer model.CloseDB(db)

	if err := cache.InitRedis(cfg); err != nil {
		log.Printf("Warning: Failed to initialize Redis: %v", err)
	}
	defer cache.CloseRedis()

	auth.InitJWT(cfg)

	h := server.Default(
		server.WithHostPorts("0.0.0.0:8888"),
	)

	registerRoutes(h, db, cfg)

	h.Spin()
}
