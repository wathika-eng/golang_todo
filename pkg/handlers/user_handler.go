package handlers

import (
	"fmt"
	"golang_todo/pkg/repository"
	"golang_todo/pkg/services"
	"golang_todo/pkg/types"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golodash/galidator"
)

// This allows all methods in UserHandler to access
//
//	database operations without directly interacting with the database.
type UserHandler struct {
	userRepo     *repository.UserRepo
	userServices *services.UserServices
}

// constructor
func NewUserHandler(userRepo *repository.UserRepo, userServices *services.UserServices) *UserHandler {
	return &UserHandler{userRepo: userRepo, userServices: userServices}
}

var (
	g         = galidator.G()
	validator = g.Validator(types.User{})
)

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
			"message": validator.DecryptErrors(err),
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

}
