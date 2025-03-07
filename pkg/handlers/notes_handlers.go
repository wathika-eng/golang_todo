package handlers

import (
	"fmt"
	"golang_todo/pkg/repository"
	notesservices "golang_todo/pkg/services/notes_services"
	"golang_todo/pkg/types"
	"net/http"
	"strconv"
	"strings"
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
			"message": notesValidator.DecryptErrors(err),
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
func (h *NotesHandler) GetNotes(c *gin.Context) {
	notes, err := h.NotesRepo.GetAllNotes()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	fmt.Println(notes)
	if len(notes) == 0 {
		c.JSON(200, gin.H{
			"error":   false,
			"message": "no todos in the database",
			"todos":   notes,
		})
		return
	}
	c.JSON(200, gin.H{
		"error":   false,
		"message": notes,
	})
}

// read
func (h *NotesHandler) GetNoteByID(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	uintID, err := strconv.Atoi(id)
	if err != nil || uintID < 0 {
		errM := fmt.Sprintf("could not convert %v to integer: %v", id, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{

			"error":   true,
			"message": errM,
		})
		return
	}
	notes, err := h.NotesRepo.GetNoteByID(uint(uintID))
	if err != nil || uintID <= 0 {
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

func (h *NotesHandler) NotesTest(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "notes API is up and running",
		"time":    time.Now().Local(),
	})
}
