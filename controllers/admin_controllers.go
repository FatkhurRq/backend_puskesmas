package controllers

import (
	"github.com/fath/puskesmas-backend/config"
	"github.com/fath/puskesmas-backend/dto"
	"github.com/fath/puskesmas-backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AdminDashboard(c *gin.Context) {

	var totalPatients int64
	var totalPolis int64

	var waitingQueues int64
	var calledQueues int64
	var completedQueues int64

	config.DB.
		Model(&models.User{}).
		Where("role = ?", "patient").
		Count(&totalPatients)

	config.DB.
		Model(&models.Poli{}).
		Count(&totalPolis)

	config.DB.
		Model(&models.Queue{}).
		Where("status = ?", "waiting").
		Count(&waitingQueues)

	config.DB.
		Model(&models.Queue{}).
		Where("status = ?", "called").
		Count(&calledQueues)

	config.DB.
		Model(&models.Queue{}).
		Where("status = ?", "completed").
		Count(&completedQueues)

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Dashboard data",
		Data: gin.H{
			"total_patients":   totalPatients,
			"total_polis":      totalPolis,
			"waiting_queues":   waitingQueues,
			"called_queues":    calledQueues,
			"completed_queues": completedQueues,
		},
	})
}

func CallQueue(c *gin.Context) {

	queueID := c.Param("id")

	var queue models.Queue

	err := config.DB.
		Preload("Registration").
		Preload("Registration.User").
		Preload("Registration.Poli").
		First(&queue, "id = ?", queueID).
		Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Queue not found",
		})
		return
	}

	if queue.Status != "waiting" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Queue cannot be called",
		})
		return
	}

	poliID := queue.Registration.PoliID

	// Cari queue yang sedang dipanggil di poli yang sama
	var activeQueue models.Queue

	err = config.DB.
		Joins("JOIN registrations ON registrations.id = queues.registration_id").
		Where("registrations.poli_id = ?", poliID).
		Where("queues.status = ?", "called").
		First(&activeQueue).Error

	// Jika ada queue aktif → selesaikan otomatis
	if err == nil {

		config.DB.Model(&models.Queue{}).
			Where("id = ?", activeQueue.ID).
			Update("status", "completed")

		config.DB.Model(&models.Registration{}).
			Where("id = ?", activeQueue.RegistrationID).
			Update("status", "completed")
	}

	// Set queue baru menjadi called
	if err := config.DB.
		Model(&models.Queue{}).
		Where("id = ?", queue.ID).
		Update("status", "called").Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Update registration
	config.DB.
		Model(&models.Registration{}).
		Where("id = ?", queue.RegistrationID).
		Update("status", "called")

	// Ambil ulang data lengkap
	config.DB.
		Preload("Registration").
		Preload("Registration.User").
		Preload("Registration.Poli").
		First(&queue, "id = ?", queue.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Queue called successfully",
		"data":    queue,
	})
}

func CompleteQueue(c *gin.Context) {

	queueID := c.Param("id")

	var queue models.Queue

	err := config.DB.
		Preload("Registration").
		Preload("Registration.User").
		Preload("Registration.Poli").
		First(&queue, "id = ?", queueID).
		Error

	if err != nil {
		c.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Message: "Queue not found",
			Data:    nil,
		})
		return
	}

	if queue.Status != "called" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Queue must be called first",
		})
		return
	}

	// Update status queue
	if err := config.DB.
		Model(&models.Queue{}).
		Where("id = ?", queue.ID).
		Update("status", "completed").Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Update status registration
	config.DB.
		Model(&models.Registration{}).
		Where("id = ?", queue.RegistrationID).
		Update("status", "completed")

	// Ambil ulang data lengkap
	config.DB.
		Preload("Registration").
		Preload("Registration.User").
		Preload("Registration.Poli").
		First(&queue, "id = ?", queue.ID)

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Queue completed successfully",
		Data:    queue,
	})
}

func GetAllQueues(c *gin.Context) {

	var queues []models.Queue

	err := config.DB.
		Preload("Registration").
		Preload("Registration.User").
		Preload("Registration.Poli").
		Find(&queues).Error

	if err != nil {

		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "Failed to retrieve queues",
			Data:    nil,
		})

		return
	}

	var response []dto.AllQueueResponse

	for _, queue := range queues {

		response = append(response, dto.AllQueueResponse{
			QueueID:     queue.ID.String(),
			QueueNumber: queue.QueueNumber,
			PatientName: queue.Registration.User.Name,
			Poli:        queue.Registration.Poli.Name,
			Status:      queue.Status,
			QueueDate:   queue.QueueDate.Format("2006-01-02"),
		})
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Queues retrieved successfully",
		Data:    response,
	})
}

func GetPatients(c *gin.Context) {

	var users []models.User

	err := config.DB.
		Where("role = ?", "patient").
		Find(&users).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	var patients []dto.PatientResponse

	for _, user := range users {
		patients = append(patients, dto.PatientResponse{
			ID:    user.ID,
			NIK:   user.NIK,
			Name:  user.Name,
			Phone: user.Phone,
		})
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Patients retrieved successfully",
		Data:    patients,
	})
}
