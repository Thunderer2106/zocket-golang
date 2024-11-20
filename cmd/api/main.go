


	package main

	import (
		"log"
		"os"
		"os/signal"
		"syscall"
	
		"zocket/internal/api"
		"zocket/redis"
		"zocket/internal/config"
		"zocket/internal/db"
		"zocket/internal/api/handlers"
		"github.com/gin-gonic/gin"
	)
	

func main() {
	cfg := config.LoadConfig()
	db.ConnectToDB(cfg)

	handlers.InitRabbitMQ()
	defer handlers.CloseRabbitMQ() 


	redis.InitRedis()
    go handlers.StartImageProcessingWorker()
	router := gin.Default()

	api.SetupRoutes(router)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Starting server on port 8080...")
		if err := router.Run(":8080"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")
}
