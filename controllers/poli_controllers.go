package controllers

import (
	"github.com/fath/puskesmas-backend/config"
	"github.com/fath/puskesmas-backend/dto"
	"github.com/fath/puskesmas-backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PoliInput struct {
	Name string `json:"name"`
}

func CreatePoli(c *gin.Context) {

	var input PoliInput

	if err := c.ShouldBindJSON(&input); err != nil {

		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid input",
			Data:    nil,
		})

		return
	}

	poli := models.Poli{
		Name: input.Name,
	}

	if err := config.DB.Create(&poli).Error; err != nil {

		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Failed to create poli",
			Data:    nil,
		})

		return
	}

	c.JSON(http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "Poli created successfully",
		Data:    poli,
	})
}

func GetPolis(c *gin.Context) {

	var polis []models.Poli

	config.DB.Find(&polis)

	c.JSON(http.StatusOK, polis)
}

func GetCurrentQueue(c *gin.Context) {

	poliID := c.Param("id")

	var queue models.Queue

	// Cari antrean yang sedang dipanggil
	err := config.DB.
		Preload("Registration").
		Preload("Registration.User").
		Preload("Registration.Poli").
		Joins("JOIN registrations ON registrations.id = queues.registration_id").
		Where("registrations.poli_id = ?", poliID).
		Where("queues.status = ?", "called").
		First(&queue).Error

	if err == nil {

		response := dto.QueueResponse{
			QueueNumber: queue.QueueNumber,
			PatientName: queue.Registration.User.Name,
			Poli:        queue.Registration.Poli.Name,
			Status:      queue.Status,
		}

		c.JSON(http.StatusOK, dto.APIResponse{
			Success: true,
			Message: "Current queue retrieved successfully",
			Data:    response,
		})

		return
	}

	// Jika belum ada yang dipanggil
	c.JSON(http.StatusOK, dto.APIResponse{
		Success: false,
		Message: "No queue currently called",
		Data:    nil,
	})

}

func GetNextQueue(c *gin.Context) {

	poliID := c.Param("id")

	var queue models.Queue

	err := config.DB.
		Preload("Registration").
		Preload("Registration.User").
		Preload("Registration.Poli").
		Joins("JOIN registrations ON registrations.id = queues.registration_id").
		Where("registrations.poli_id = ?", poliID).
		Where("queues.status = ?", "waiting").
		Order("queues.queue_number ASC").
		First(&queue).Error

	if err != nil {
		c.JSON(http.StatusOK, dto.APIResponse{
			Success: false,
			Message: "No waiting queue",
			Data:    nil,
		})
		return
	}

	response := dto.QueueResponse{
		QueueNumber: queue.QueueNumber,
		PatientName: queue.Registration.User.Name,
		Poli:        queue.Registration.Poli.Name,
		Status:      queue.Status,
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Next queue retrieved successfully",
		Data:    response,
	})
}
