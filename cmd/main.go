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
	r.Use(cors.Default())

	routes.SetupRoutes(r)

	r.Run(":8080")
}
