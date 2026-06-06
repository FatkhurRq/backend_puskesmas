package routes

import (
	"github.com/fath/puskesmas-backend/controllers"
	"github.com/fath/puskesmas-backend/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	api := r.Group("/api")

	{
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)
		api.GET("/me", middleware.JWTAuthMiddleware(), controllers.Me)
		api.POST("/polis", controllers.CreatePoli)
		api.GET("/polis", controllers.GetPolis)
		api.POST("/registrations", middleware.JWTAuthMiddleware(), controllers.CreateRegistration)
		api.GET("/registrations/me", middleware.JWTAuthMiddleware(), controllers.GetMyRegistrations)
		api.GET("/queue/me", middleware.JWTAuthMiddleware(), controllers.GetMyQueue)
		api.GET("/admin/dashboard", middleware.JWTAuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.AdminDashboard)
		api.GET("/admin/patients", middleware.JWTAuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.GetPatients)
		api.PATCH("/queues/:id/call", middleware.JWTAuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.CallQueue)
		api.PATCH("/queues/:id/complete", middleware.JWTAuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.CompleteQueue)
		api.GET("/queues", middleware.JWTAuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.GetAllQueues)
		api.GET("/polis/:id/queue", controllers.GetCurrentQueue)
	}
}
