package redis

import (
	"context"
	"log"
	"os"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"strconv" 
)

var client *redis.Client
var ctx = context.Background()

func InitRedis() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables.")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := os.Getenv("REDIS_DB")

	db, err := strconv.Atoi(redisDB)
	if err != nil {
		log.Fatalf("Error converting REDIS_DB to integer: %v", err)
	}

	client = redis.NewClient(&redis.Options{
		Addr:     redisAddr,     
		Password: redisPassword, 
		DB:       db,            
	})
	_, err = client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	log.Println("Successfully connected to Redis")
}

func GetClient() *redis.Client {
	return client
}
// package redis

// import (
// 	"context"
// 	"log"
// 	"os"
// 	"strconv"

// 	"github.com/go-redis/redis/v8"
// )

// var client *redis.Client

// func InitRedis() {
// 	addr := os.Getenv("REDIS_ADDR")
// 	if addr == "" {
// 		addr = "localhost:6379" 
// 	}

// 	// password := os.Getenv("REDIS_PASSWORD")
// 	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
// 	if err != nil {
// 		db = 0
// 	}

// 	client = redis.NewClient(&redis.Options{
// 		Addr:     addr,
// 		Password: "foobared", 
// 		DB:       db,       
// 	})


// 	_, err = client.Ping(context.Background()).Result()
// 	if err != nil {
// 		log.Fatalf("Unable to connect to Redis: %v", err)
// 	}
// }

// func GetClient() *redis.Client {
// 	return client
// }
