package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"receipt-processor-challenge/handlers"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/receipts/process", handlers.ProcessReceipt)
	r.GET("/receipts/:id/points", handlers.GetPoints)
	return r
}

// This function tests processReceipt with a valid receipt.
// It checks to see the response returns a StatusCode of 200 and JSON with the id field.
func TestProcessReceipt_ValidInput(t *testing.T) {
	router := setupRouter()

	payload := `{
		"retailer": "Target",
		"purchaseDate": "2022-01-01",
		"purchaseTime": "13:01",
		"items": [
			{"shortDescription": "Item 1", "price": "6.49"}
		],
		"total": "6.49"
	}`

	req, _ := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "id")
}

// This function tests processReceipt with an invalid receipt.
// It checks to see if the response returns a StatusCode of 400 and JSON with the error field.
func TestProcessReceipt_InvalidInput(t *testing.T) {
	router := setupRouter()

	payload := `{
		"retailer": "Target",
		"items": [
			{"shortDescription": "Item 1"}
		]
	}`

	req, _ := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "error")
}

// This functin tests GetPoints with a valid ID.
// It checks to see if the response returns a StatusCode of 200 and JSON with the points field.
func TestGetPoints_ValidID(t *testing.T) {
	router := setupRouter()

	handlers.ReceiptStore["test-id"] = 50

	req, _ := http.NewRequest("GET", "/receipts/test-id/points", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "points")
	assert.Contains(t, resp.Body.String(), "50")
}

// This function tests GetPoints with an invalid ID.
// It checks to see if the response returns a StatusCode of 404 and JSON with the error field.
func TestGetPoints_InvalidID(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/receipts/non-existent-id/points", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
	assert.Contains(t, resp.Body.String(), "error")
}

// This function tests both processReceipt and GetPoints together.
// It checks to see if the calculations of the points are correct.
func TestProcessAndGetPoints(t *testing.T) {
	router := setupRouter()

	payload := `{
		"retailer": "Target",
		"purchaseDate": "2022-01-01",
		"purchaseTime": "13:01",
		"items": [
			{ "shortDescription": "Mountain Dew 12PK", "price": "6.49" },
			{ "shortDescription": "Emils Cheese Pizza", "price": "12.25" },
			{ "shortDescription": "Knorr Creamy Chicken", "price": "1.26" },
			{ "shortDescription": "Doritos Nacho Cheese", "price": "3.35" },
			{ "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ", "price": "12.00" }
		],
		"total": "35.35"
	}`

	req, _ := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var result map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)
	id, exists := result["id"]
	assert.True(t, exists, "Response should contain an id")

	getReq, _ := http.NewRequest("GET", "/receipts/"+id+"/points", nil)
	getResp := httptest.NewRecorder()
	router.ServeHTTP(getResp, getReq)

	assert.Equal(t, http.StatusOK, getResp.Code)
	var pointsResult map[string]int
	err = json.Unmarshal(getResp.Body.Bytes(), &pointsResult)
	assert.NoError(t, err)
	points, exists := pointsResult["points"]
	assert.True(t, exists, "Response should contain points")

	expectedPoints := 28
	assert.Equal(t, expectedPoints, points, "Points should match the expected calculation")
}
