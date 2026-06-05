package controllers

import (
	"github.com/fath/puskesmas-backend/config"
	"github.com/fath/puskesmas-backend/dto"
	"github.com/fath/puskesmas-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type RegistrationInput struct {
	PoliID string `json:"poli_id"`
}

func CreateRegistration(c *gin.Context) {

	var input RegistrationInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	userIDValue, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "User not found in token",
			Data:    nil,
		})
		return
	}

	userIDString, ok := userIDValue.(string)

	if !ok {
		c.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "Invalid user id",
			Data:    nil,
		})
		return
	}

	userUUID, err := uuid.Parse(userIDString)

	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid user UUID",
			Data:    nil,
		})
		return
	}

	poliUUID, err := uuid.Parse(input.PoliID)

	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid poli UUID",
			Data:    nil,
		})
		return
	}

	// Validasi poli
	var poli models.Poli

	err = config.DB.
		First(&poli, "id = ?", poliUUID).
		Error

	if err != nil {
		c.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Message: "Poli not found",
			Data:    nil,
		})
		return
	}

	// Cegah daftar dua kali pada poli yang sama di hari yang sama
	var existingRegistration models.Registration

	today := time.Now().Format("2006-01-02")

	err = config.DB.
		Where("user_id = ?", userUUID).
		Where("poli_id = ?", poliUUID).
		Where("DATE(registration_date) = ?", today).
		First(&existingRegistration).
		Error

	if err == nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "You already registered for this poli today",
			Data:    nil,
		})
		return
	}

	// Buat registration
	registration := models.Registration{
		UserID:           userUUID,
		PoliID:           poliUUID,
		RegistrationDate: time.Now(),
		Status:           "waiting",
	}

	if err := config.DB.Create(&registration).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	// Cari antrean terakhir berdasarkan poli
	var lastQueue models.Queue

	config.DB.
		Joins("JOIN registrations ON registrations.id = queues.registration_id").
		Where("registrations.poli_id = ?", poliUUID).
		Order("queues.queue_number DESC").
		First(&lastQueue)

	nextQueueNumber := lastQueue.QueueNumber + 1

	// Buat queue
	queue := models.Queue{
		RegistrationID: registration.ID,
		QueueNumber:    nextQueueNumber,
		QueueDate:      time.Now(),
		Status:         "waiting",
	}

	if err := config.DB.Create(&queue).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	// Ambil queue lengkap
	var resultQueue models.Queue

	err = config.DB.
		Preload("Registration").
		Preload("Registration.User").
		Preload("Registration.Poli").
		First(&resultQueue, "id = ?", queue.ID).
		Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "Registration created successfully",
		Data:    resultQueue,
	})
}

func GetMyRegistrations(c *gin.Context) {

	userIDValue, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "User not found in token",
			Data:    nil,
		})
		return
	}

	var registrations []models.Registration

	config.DB.
		Preload("Poli").
		Where("user_id = ?", userIDValue).
		Order("registration_date DESC").
		Find(&registrations)

	c.JSON(http.StatusOK, registrations)
}

func GetMyQueue(c *gin.Context) {

	userIDValue, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "User not found in token",
			Data:    nil,
		})
		return
	}

	var registration models.Registration

	err := config.DB.
		Preload("User").
		Preload("Poli").
		Where("user_id = ?", userIDValue).
		Order("registration_date DESC").
		First(&registration).Error

	if err != nil {
		c.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Message: "Registration not found",
			Data:    nil,
		})
		return
	}

	var queue models.Queue

	err = config.DB.
		Where("registration_id = ?", registration.ID).
		First(&queue).Error

	if err != nil {
		c.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Message: "Queue not found",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Queue retrieved successfully",
		Data: gin.H{
			"queue_id":          queue.ID,
			"queue_number":      queue.QueueNumber,
			"status":            queue.Status,
			"patient_name":      registration.User.Name,
			"poli":              registration.Poli.Name,
			"registration_date": registration.RegistrationDate,
		},
	})
}
