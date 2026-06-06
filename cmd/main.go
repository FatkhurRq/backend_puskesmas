package main

import (
	"github.com/fath/puskesmas-backend/config"
	"github.com/fath/puskesmas-backend/models"
	"github.com/fath/puskesmas-backend/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	config.ConnectDatabase()

	config.DB.AutoMigrate(
		&models.User{},
		&models.Poli{},
		&models.Registration{},
		&models.Queue{},
	)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	routes.SetupRoutes(r)

	r.Run(":8081")
}
