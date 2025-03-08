package handlers

import (
	"fmt"
	"golang_todo/pkg/repository"
	"golang_todo/pkg/services"
	"golang_todo/pkg/types"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// This allows all methods in UserHandler to access
//
//	database operations without directly interacting with the database.
type UserHandler struct {
	userRepo     *repository.UserRepo
	userServices services.Auth
}

// constructor
func NewUserHandler(userRepo *repository.UserRepo, userServices services.Auth) *UserHandler {
	return &UserHandler{
		userRepo:     userRepo,
		userServices: userServices,
	}
}

// var (
// 	g              = galidator.G()
// 	userValidator  = g.Validator(types.User{})
// 	notesValidator = g.Validator(types.Note{})
// )

func (h *UserHandler) SignUp(c *gin.Context) {
	//c.Header("Cache-Control", "no-store, no-cache, must-revalidate, private")
	var user types.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	hashedPassword, err := h.userServices.HashPassword(user.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	user.Password = hashedPassword
	// _, _ = h.userServices.SendEmail(user.Email)
	err = h.userRepo.InsertUser(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(201, gin.H{
		"error":   false,
		"message": fmt.Sprintf("user: %v created successfully", user.Email),
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	//c.Header("Cache-Control", "no-store, no-cache, must-revalidate, private")
	var user types.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	userFound, err := h.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	err = h.userServices.CheckPassword(userFound.Password, user.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Wrong password!",
		})
		return
	}
	access_token, refresh_token, err := h.userServices.GenerateToken(userFound.ID, userFound.Email, user.Role, false)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Wrong password!",
		})
		return
	}

	c.JSON(200, gin.H{
		"error":         false,
		"message":       "Access granted",
		"access_token":  access_token,
		"refresh_token": refresh_token,
	})
}

func (h *UserHandler) RefreshAccess(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil || strings.TrimSpace(req.RefreshToken) == "" {
		c.AbortWithStatusJSON(400, gin.H{
			"error":   true,
			"message": "refresh token not found",
		})
		return
	}
	validToken, err := h.userServices.ValidateToken(req.RefreshToken, true)
	if err != nil || !validToken.Valid || validToken == nil {
		// log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{
			"error":   true,
			"message": "refresh token is invalid",
		})
		return
	}
	claims, ok := validToken.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatusJSON(400, gin.H{
			"error":   true,
			"message": "invalid token data",
		})
		return
	}
	userID, email, role := claims["user_id"].(uuid.UUID), claims["sub"].(string), claims["role"].(string)
	newAccessToken, newRefreshToken, err := h.userServices.GenerateToken(userID, email, role, true)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error":   true,
			"message": "failed to generate new tokens",
		})
		return
	}
	c.JSON(200, gin.H{
		"error":         false,
		"message":       "Access granted",
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

func (r *UserHandler) UserProfile(c *gin.Context) {
	userEmail, exists := c.Get("user_email")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "unauthorized",
		})
		return
	}
	userData, err := r.userRepo.GetUserByEmail(userEmail.(string))
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error":   true,
			"message": "failed to fetch user details",
		})
		return
	}
	userData.Password = ""
	c.JSON(200, gin.H{
		"user": userData,
	})
}

func (r *UserHandler) UserTest(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "API is up and running",
		"time":    time.Now().Local(),
	})
}
