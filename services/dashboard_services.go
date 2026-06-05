package services

import (
	"github.com/fath/puskesmas-backend/config"
	"github.com/fath/puskesmas-backend/dto"
	"github.com/fath/puskesmas-backend/models"
)

func GetDashboardData() (dto.DashboardResponse, error) {
	var dashboard dto.DashboardResponse

	// Total users
	if err := config.DB.Model(&models.User{}).
		Count(&dashboard.TotalUsers).Error; err != nil {
		return dashboard, err
	}

	// Total poli
	if err := config.DB.Model(&models.Poli{}).
		Count(&dashboard.TotalPoli).Error; err != nil {
		return dashboard, err
	}

	// Total waiting
	if err := config.DB.Model(&models.Queue{}).
		Where("status = ?", "waiting").
		Count(&dashboard.TotalWaiting).Error; err != nil {
		return dashboard, err
	}

	// Total called
	if err := config.DB.Model(&models.Queue{}).
		Where("status = ?", "called").
		Count(&dashboard.TotalCalled).Error; err != nil {
		return dashboard, err
	}

	// Total completed
	if err := config.DB.Model(&models.Queue{}).
		Where("status = ?", "completed").
		Count(&dashboard.TotalCompleted).Error; err != nil {
		return dashboard, err
	}

	return dashboard, nil
}
