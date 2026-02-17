package handlers

import (
	"locker-storage-api/config"
	"locker-storage-api/internal/models"
	"locker-storage-api/internal/utils"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phoneNumber"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	// Check if user exists
	var existingUser models.User
	if err := config.GetDB().Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		utils.ErrorResponse(c, 400, "Email already registered")
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to hash password")
		return
	}

	// Create user
	user := models.User{
		Email:       req.Email,
		Password:    hashedPassword,
		Name:        req.Name,
		PhoneNumber: &req.PhoneNumber,
		IsActive:    true,
	}

	if err := config.GetDB().Create(&user).Error; err != nil {
		utils.ErrorResponse(c, 500, "Failed to create user")
		return
	}

	// Generate JWT
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to generate token")
		return
	}

	utils.SuccessResponse(c, 201, "User registered successfully", gin.H{
		"user":  user,
		"token": token,
	})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	// Find user
	var user models.User
	if err := config.GetDB().Where("email = ?", req.Email).First(&user).Error; err != nil {
		utils.ErrorResponse(c, 401, "Invalid credentials")
		return
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		utils.ErrorResponse(c, 401, "Invalid credentials")
		return
	}

	// Generate JWT
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to generate token")
		return
	}

	utils.SuccessResponse(c, 200, "Login successful", gin.H{
		"user":  user,
		"token": token,
	})
}

func GetProfile(c *gin.Context) {
	user, _ := c.Get("user")
	utils.SuccessResponse(c, 200, "Profile retrieved", user)
}
