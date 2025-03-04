package handlers

import (
	"fmt"
	"golang_todo/pkg/repository"
	"golang_todo/pkg/services"

	"github.com/gin-gonic/gin"
)

type NotesHandler struct{
	NotesRepo *repository.NotesRepository
	NotesServices *services.NotesServices
}

func NewNotesHandler(repository.NotesRepository, services.NotesServices)  {

}

// create
func CreateNotes(c *gin.Context) {

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
