package handlers

import (
	"fmt"
	"golang_todo/pkg/repository"
	"golang_todo/pkg/response"
	notesservices "golang_todo/pkg/services/notes_services"
	redisservices "golang_todo/pkg/services/redis"
	"golang_todo/pkg/types"
	"log"
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
	response      response.ResponseInterface
}

func NewNotesHandler(notesRepo *repository.NotesRepository, notesServices *notesservices.NotesServices,
	redisServices redisservices.Redis, response response.ResponseInterface) *NotesHandler {
	return &NotesHandler{
		NotesRepo:     notesRepo,
		NotesServices: notesServices,
		redisServices: redisServices,
		response:      response,
	}
}

// CreateNotes create
func (h *NotesHandler) CreateNotes(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.response.SendError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var notes types.Note
	err := c.ShouldBindBodyWithJSON(&notes)
	if err != nil {
		h.response.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	notes.UserID = userID.(uuid.UUID)
	err = h.NotesRepo.InsertNotes(notes)
	if err != nil {
		h.response.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.redisServices.DeleteCache(userID.(uuid.UUID))
	if err != nil {
		h.response.SendError(c, http.StatusInternalServerError, err.Error())
	}
	h.response.Success(c, 201, "note created successfuly")
}

// read
func (h *NotesHandler) GetNotes(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.response.SendError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	notes, err := h.redisServices.FetchFromCache(userID.(uuid.UUID))
	if err == nil || len(notes) != 0 {
		log.Println("fetched from cache")
		c.Header("X-Cache-Status", "HIT")
		h.response.Success(c, 200, notes)
		return
	}
	notes, err = h.NotesRepo.GetAllNotes(userID.(uuid.UUID))
	if err != nil {
		h.response.SendError(c, http.StatusBadRequest, err.Error())
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
	err = h.redisServices.CacheTodo(notes, userID.(uuid.UUID))
	// if err != nil {
	// 	log.Printf("error while trying to cache: %v", err.Error())
	// }
	log.Println("fetched from db")
	c.Header("X-Cache-Status", "MISS")
	h.response.Success(c, 200, notes)
}

// read
func (h *NotesHandler) GetNoteByID(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	uintID, err := uuid.Parse(id)
	if err != nil {
		errM := fmt.Sprintf("could not convert %v to uuid: %v", id, err.Error())
		h.response.SendError(c, http.StatusBadRequest, errM)
		return
	}
	note, err := h.NotesRepo.GetNoteByID(uintID)
	if err != nil {
		h.response.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.redisServices.CacheTodo(note, uintID)
	if err != nil {
		h.response.SendError(c, http.StatusInternalServerError, err.Error())
	}
	h.response.Success(c, 200, note)
}

// update
func (h *NotesHandler) UpdateNotes(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	userID, _ := c.Get("user_id")
	uintID, err := uuid.Parse(id)
	if err != nil {
		errM := fmt.Sprintf("could not convert %v to uuid: %v", id, err.Error())
		h.response.SendError(c, http.StatusBadRequest, errM)
		return
	}
	var toUpdateFields map[string]interface{}
	err = c.ShouldBindJSON(&toUpdateFields)
	if err != nil {
		h.response.SendError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if body, ok := toUpdateFields["body"].(string); ok {
		cleanedBody := strings.TrimSpace(body)
		if cleanedBody == "" {
			h.response.SendError(c, http.StatusBadRequest, "body cannot be empty or blank")
			return
		}
		toUpdateFields["body"] = cleanedBody
	}

	newTodo, err := h.NotesRepo.UpdateWithID(uintID, toUpdateFields)
	if err != nil || newTodo == nil {
		h.response.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	h.redisServices.DeleteCache(userID.(uuid.UUID))
	c.JSON(200, gin.H{
		"error":   false,
		"message": "updated todo successfully",
		"todo":    newTodo,
	})
}

// delete
func (h *NotesHandler) DeleteNotes(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	userID, _ := c.Get("user_id")
	uintID, err := uuid.Parse(id)
	if err != nil {
		errM := fmt.Sprintf("could not convert %v to uuid: %v", id, err.Error())
		h.response.SendError(c, 400, errM)
		return
	}
	ok, err := h.NotesRepo.DeleteWithID(uintID)
	if err != nil || !ok {
		h.response.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	h.redisServices.DeleteCache(userID.(uuid.UUID))
	c.JSON(200, gin.H{
		"error":   false,
		"message": "deleted todo successfully",
		"todo":    "[]",
	})
}

func (h *NotesHandler) RecentDeletions(c *gin.Context) {
	userID, _ := c.Get("user_id")
	notes, err := h.NotesRepo.SoftDelete(userID.(uuid.UUID))
	if err != nil {
		h.response.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	if len(notes) <= 0 {
		c.JSON(200, gin.H{
			"error":   false,
			"message": "no recently deleted todos",
		})
		return
	}
	h.response.Success(c, 200, notes)
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
	h.response.Success(c, 200, time.Now().Local())
}
