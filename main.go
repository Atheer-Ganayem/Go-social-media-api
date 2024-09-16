package main

import (
	"os"

	db "github.com/Atheer-Ganayem/Go-social-media-api/DB"
	"github.com/Atheer-Ganayem/Go-social-media-api/routes"
	"github.com/Atheer-Ganayem/Go-social-media-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	db.Init()
	defer db.Disconnect()
	utils.InitAWS()

	server := gin.Default()
	routes.Register(server)
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "127.0.0.1:8080"
	} else {
		port = ":" + port
	}

	server.Run(port)
}