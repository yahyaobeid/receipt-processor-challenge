package main

import (
	"log"
	"receipt-processor-challenge/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// I want to use a router for organization and best practices
	router := gin.Default()

	router.POST("/receipts/process", handlers.ProcessReceipt)
	router.GET("/receipts/:id/points", handlers.GetPoints)

	log.Println("Server is running on Port 8000...")
	if err := router.Run(":8000"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
