package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

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
