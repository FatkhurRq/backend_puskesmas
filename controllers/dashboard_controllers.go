package controllers

import (
	"github.com/fath/puskesmas-backend/dto"
	"github.com/fath/puskesmas-backend/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetDashboard(c *gin.Context) {
	dashboard, err := services.GetDashboardData()

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "Failed to get dashboard data",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Dashboard retrieved successfully",
		Data:    dashboard,
	})
}
