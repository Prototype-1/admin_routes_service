package main

import (
	"github.com/joho/godotenv"
    "log"
	"github.com/Prototype-1/admin_routes_service/internal/server"
)

func init() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
}

func main() {
	server.StartServer()
}
