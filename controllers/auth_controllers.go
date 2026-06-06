package controllers

import (
	"net/http"

	"github.com/fath/puskesmas-backend/config"
	"github.com/fath/puskesmas-backend/dto"
	"github.com/fath/puskesmas-backend/models"
	"github.com/fath/puskesmas-backend/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	NIK      string `json:"nik"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type LoginInput struct {
	NIK      string `json:"nik"`
	Password string `json:"password"`
}

func Register(c *gin.Context) {

	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {

		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})

		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "failed to hash password",
			Data:    nil,
		})
		return
	}

	user := models.User{
		NIK:      input.NIK,
		Name:     input.Name,
		Phone:    input.Phone,
		Password: string(hashedPassword),
		Role:     "patient",
	}

	if err := config.DB.Create(&user).Error; err != nil {

		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})

		return
	}

	c.JSON(http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "Register success",
		Data:    nil,
	})
}

func Login(c *gin.Context) {

	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	var user models.User

	result := config.DB.
		Where("nik = ?", input.NIK).
		First(&user)

	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "NIK or Password incorrect",
			Data:    nil,
		})
		return
	}

	err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(input.Password),
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "NIK or Password incorrect",
			Data:    nil,
		})
		return
	}

	token, err := utils.GenerateToken(
		user.ID.String(),
		user.Role,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "Failed to generate token",
			Data:    nil,
		})
		return
	}

	userResponse := dto.PatientResponse{
		ID:    user.ID,
		NIK:   user.NIK,
		Name:  user.Name,
		Phone: user.Phone,
		Role:  user.Role,
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Login success",
		Data: gin.H{
			"token": token,
			"user":  userResponse,
		},
	})
}

func Me(c *gin.Context) {

	userID, _ := c.Get("user_id")

	var user models.User
	result := config.DB.Where("id = ?", userID).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "User not found",
			Data:    nil,
		})
		return
	}

	userResponse := dto.PatientResponse{
		ID:    user.ID,
		NIK:   user.NIK,
		Name:  user.Name,
		Phone: user.Phone,
		Role:  user.Role,
	}

	c.JSON(200, dto.APIResponse{
		Success: true,
		Message: "User data retrieved successfully",
		Data:    gin.H{"user": userResponse},
	})
}
