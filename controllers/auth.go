package controllers

import (
	"ecommerce-backend/config"
	"ecommerce-backend/models"
	"ecommerce-backend/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	FullName string `json:"fullName" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses password"})
		return
	}

	query := `INSERT INTO users (full_name, email, password_hash) VALUES ($1, $2, $3) RETURNING id, role, created_at`
	var user models.User

	err = config.DB.QueryRowx(query, input.FullName, input.Email, string(hashedPassword)).StructScan(&user)
	if err != nil {
		fmt.Println("DB Error Detail:", err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Email mungkin sudah terdaftar atau terjadi kesalahan server"})
		return
	}

	user.FullName = input.FullName
	user.Email = input.Email

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registrasi berhasil",
		"data":    user,
	})
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Email atau password salah", err)
		return
	}

	var user models.User
	query := `SELECT id, email, password_hash, role FROM users WHERE email = $1`

	err := config.DB.Get(&user, query, input.Email)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Email atau password salah", err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Email atau password salah", err)
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal membuat sesi login", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login berhasil", gin.H{
		"token": token,
	})
}
