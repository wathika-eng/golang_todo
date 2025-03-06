package handlers

import (
	"fmt"
	"golang_todo/pkg/repository"
	notesservices "golang_todo/pkg/services/notes_services"
	"golang_todo/pkg/types"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type NotesHandler struct {
	NotesRepo     *repository.NotesRepository
	NotesServices *notesservices.NotesServices
}

func NewNotesHandler(notesRepo *repository.NotesRepository, notesServices *notesservices.NotesServices) *NotesHandler {
	return &NotesHandler{
		NotesRepo:     notesRepo,
		NotesServices: notesServices,
	}
}

// create
func (h *NotesHandler) CreateNotes(c *gin.Context) {
	var notes types.Note
	err := c.ShouldBindBodyWithJSON(&notes)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
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
func GetNotes(c *gin.Context) {

}

// read
func GetNoteByID(c *gin.Context) {
	id := c.Param("id")
	resp := fmt.Sprintf("id requested: %v", id)
	c.JSON(200, gin.H{
		"message": resp,
	})
}

// update
func UpdateNotes(c *gin.Context) {

}

// delete
func DeleteNotes(c *gin.Context) {

}

func GetUserDetails(c *gin.Context) {

}

func (h *NotesHandler) NotesTest(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "notes API is up and running",
		"time":    time.Now().Local(),
	})
}
