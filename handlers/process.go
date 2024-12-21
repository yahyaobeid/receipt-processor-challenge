package handlers

import (
	"net/http"
	"receipt-processor-challenge/models"
	"receipt-processor-challenge/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var ReceiptStore = make(map[string]int)

// This function ensures the date string is a valid date
func isValidDate(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

// This function ensures the time string is a valid time
func isValidTime(timeStr string) bool {
	_, err := time.Parse("15:04", timeStr)
	return err == nil
}

// This function ensures the decimal string is a valid decimal number
func isValidDecimal(decimal string) bool {
	_, err := strconv.ParseFloat(decimal, 64)
	return err == nil
}

// This function ensures all items in the receipt have a non-empty short description and a valid price
func allItemsValid(items []models.Item) bool {
	for _, item := range items {
		if item.ShortDescription == "" || !isValidDecimal(item.Price) {
			return false
		}
	}
	return true
}

func ProcessReceipt(c *gin.Context) {
	var receipt models.Receipt
	// I ran a validation check on the components of the receipt. Part of this check was ensuring
	// that the string fields had the format that was shown in the ReadME.
	if err := c.ShouldBindJSON(&receipt); err != nil ||
		receipt.Retailer == "" ||
		len(receipt.Items) == 0 ||
		!isValidDate(receipt.PurchaseDate) ||
		!isValidTime(receipt.PurchaseTime) ||
		!isValidDecimal(receipt.Total) ||
		!allItemsValid(receipt.Items) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The receipt is invalid."})
		return
	}

	points := models.CalculatePoints(receipt)

	id := utils.GenerateID()
	ReceiptStore[id] = points

	c.JSON(http.StatusOK, gin.H{"id": id})
}
