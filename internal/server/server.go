package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "admin_routes_service is running!"})
	})

	fmt.Println("Starting admin_routes_service on port 8081...")
	if err := r.Run(":8081"); err != nil {
		log.Fatal(err)
	}
}
