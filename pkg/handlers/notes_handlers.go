package handlers

import (
	"fmt"
	"golang_todo/pkg/repository"
	notesservices "golang_todo/pkg/services/notes_services"
	redisservices "golang_todo/pkg/services/redis"
	"golang_todo/pkg/types"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotesHandler struct {
	NotesRepo     *repository.NotesRepository
	NotesServices *notesservices.NotesServices
	redisServices redisservices.Redis
}

func NewNotesHandler(notesRepo *repository.NotesRepository, notesServices *notesservices.NotesServices,
	redisServices redisservices.Redis) *NotesHandler {
	return &NotesHandler{
		NotesRepo:     notesRepo,
		NotesServices: notesServices,
		redisServices: redisServices,
	}
}

// create
func (h *NotesHandler) CreateNotes(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "unauthorized",
		})
		return
	}

	var notes types.Note
	err := c.ShouldBindBodyWithJSON(&notes)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": notesValidator.DecryptErrors(err),
		})
		return
	}
	notes.UserID = userID.(uuid.UUID)
	fmt.Println(notes.UserID)
	err = h.NotesRepo.InsertNotes(notes)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(201, gin.H{
		"error":   false,
		"message": "Notes created successfulyy",
	})
}

// read
func (h *NotesHandler) GetNotes(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "unauthorized",
		})
		return
	}
	notes, err := h.NotesRepo.GetAllNotes(userID.(uuid.UUID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	if len(notes) <= 0 {
		c.JSON(200, gin.H{
			"error":   false,
			"message": "no todos in the database",
			"todos":   notes,
		})
		return
	}
	c.JSON(200, gin.H{
		"error": false,
		"todos": notes,
	})
}

// read
func (h *NotesHandler) GetNoteByID(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	uintID, err := uuid.Parse(id)
	if err != nil {
		errM := fmt.Sprintf("could not convert %v to integer: %v", id, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{

			"error":   true,
			"message": errM,
		})
		return
	}
	notes, err := h.NotesRepo.GetNoteByID(uintID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error":   false,
		"message": notes,
	})
}

// update
func (h *NotesHandler) UpdateNotes(c *gin.Context) {

}

// delete
func (h *NotesHandler) DeleteNotes(c *gin.Context) {

}

func GetUserDetails(c *gin.Context) {

}

func (h *NotesHandler) Logout(c *gin.Context) {
	requiredKeys := []string{"user_id", "exp_time", "user_token"}
	values := make(map[string]interface{})

	for _, key := range requiredKeys {
		val, exists := c.Get(key)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": key + " not found",
			})
			return
		}
		values[key] = val
	}

	token, _ := values["user_token"].(string)
	expTimeDuration, ok := values["exp_time"].(time.Duration)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "invalid exp_time format",
		})
		return
	}
	err := h.redisServices.BlackListToken(token, expTimeDuration)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "could not blacklist token",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "logged out successfully",
	})
}

func (h *NotesHandler) NotesTest(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "notes API is up and running",
		"time":    time.Now().Local(),
	})
}
