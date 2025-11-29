package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func main() {
	ctx := context.Background()
	err := rdb.Set(ctx, "goredistest", "hello-redis", 0).Err()
	if err != nil {
		panic(err)
	}

	val, _ := rdb.Get(ctx, "goredistest").Result()
	fmt.Println(val)

	result, err := rdb.Do(ctx, "get", "goredistest").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(result.(string))
}
