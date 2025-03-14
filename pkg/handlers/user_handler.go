package handlers

import (
	"fmt"
	"golang_todo/pkg/repository"
	"golang_todo/pkg/response"
	"golang_todo/pkg/services"
	"golang_todo/pkg/types"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// UserHandler This allows all methods in UserHandler to access
//
//	database operations without directly interacting with the database.
type UserHandler struct {
	userRepo     *repository.UserRepo
	userServices services.Auth
	response     response.ResponseInterface
}

// NewUserHandler constructor
func NewUserHandler(userRepo *repository.UserRepo,
	userServices services.Auth, response response.ResponseInterface) *UserHandler {
	return &UserHandler{
		userRepo:     userRepo,
		userServices: userServices,
		response:     response,
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
	err := c.ShouldBindBodyWithJSON(&user)
	if err != nil {
		h.response.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	hashedPassword, err := h.userServices.HashPassword(user.Password)
	if err != nil {
		h.response.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	user.Password = hashedPassword
	// _, _ = h.userServices.SendEmail(user.Email)
	err = h.userRepo.InsertUser(user)
	if err != nil {
		h.response.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	// use strings builder next
	h.response.Success(c, http.StatusCreated, fmt.Sprintf("user %v created successfuly", user.Email))
}

func (h *UserHandler) Login(c *gin.Context) {
	//c.Header("Cache-Control", "no-store, no-cache, must-revalidate, private")
	var user types.User
	err := c.ShouldBindJSON(&user)
	// don't check for blank passwords or greater than 100
	if err != nil {
		h.response.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	if len(strings.TrimSpace(user.Password)) <= 0 || len(user.Password) > 100 {
		h.response.SendError(c, 400, "error with password length")
		return
	}
	userFound, err := h.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		h.response.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.userServices.CheckPassword(userFound.Password, user.Password)
	if err != nil {
		h.response.SendError(c, http.StatusBadRequest, "wrong password")
		return
	}
	accessToken, _, err := h.userServices.GenerateToken(userFound.ID, userFound.Email, user.Role, false)
	if err != nil {
		h.response.SendError(c, http.StatusBadRequest, "wrong password")
		return
	}

	c.JSON(200, gin.H{
		"error":         false,
		"message":       "Access granted",
		"access_token":  accessToken,
		"refresh_token": "",
	})
}

func (h *UserHandler) RefreshAccess(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil || strings.TrimSpace(req.RefreshToken) == "" {
		h.response.SendError(c, http.StatusBadRequest, "refresh token not found")
		return
	}
	validToken, err := h.userServices.ValidateToken(req.RefreshToken, true)
	if err != nil || !validToken.Valid {
		// log.Println(err)
		h.response.SendError(c, http.StatusBadRequest, "refresh token is invalid")
		return
	}
	claims, ok := validToken.Claims.(jwt.MapClaims)
	if !ok {
		h.response.SendError(c, http.StatusBadRequest, "invalid token")
		return
	}
	userID, email, role := claims["user_id"].(uuid.UUID), claims["sub"].(string), claims["role"].(string)
	newAccessToken, newRefreshToken, err := h.userServices.GenerateToken(userID, email, role, true)
	if err != nil {
		h.response.SendError(c, http.StatusBadRequest, "failed to generate new token")
		return
	}
	c.JSON(200, gin.H{
		"error":         false,
		"message":       "Access granted",
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

func (h *UserHandler) UserProfile(c *gin.Context) {
	userEmail, exists := c.Get("user_email")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "unauthorized",
		})
		return
	}
	userData, err := h.userRepo.GetUserByEmail(userEmail.(string))
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

func (h *UserHandler) UserTest(c *gin.Context) {
	h.response.Success(c, 200, time.Now().Local())
}
