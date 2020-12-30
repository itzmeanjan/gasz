package pubsub

import (
	"context"
	"fmt"
	"gas-price-notifier/config"
	"log"
	"strconv"

	"github.com/go-redis/redis/v8"
)

// Connect - Connects to Redis instance
func Connect() *redis.Client {

	db, err := strconv.Atoi(config.Get("RedisDatabase"))
	if err != nil {
		log.Fatalf("[!] Bad Redis database : %s\n", err.Error())
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Get("RedisHost"), config.Get("RedisPort")),
		Password: config.Get("RedisPassword"),
		DB:       db,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("[!] Failed to connect to Redis : %s\n", err.Error())
	}

	return client

}
