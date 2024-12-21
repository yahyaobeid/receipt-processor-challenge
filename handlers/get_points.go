package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPoints(c *gin.Context) {
	id := c.Param("id")

	points, exists := ReceiptStore[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "No receipt found for that ID."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"points": points})
}
