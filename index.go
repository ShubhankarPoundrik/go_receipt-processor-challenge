// This program implements a simple HTTP server that processes receipts and calculates points based on various conditions.

package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Item represents an item in a receipt.
type Item struct {
	ShortDescription string  `json:"shortDescription"`
	Price            string `json:"price"`
}

// savedReceipts stores the calculated points for each receipt.
var savedReceipts = make(map[string]int)

func main() {
	r := mux.NewRouter()

	// Use the jsonMiddleware to set the content type to JSON for all responses.
	r.Use(jsonMiddleware)

	// Define routes for receipt processing and points retrieval.
	r.HandleFunc("/receipts/process", ProcessHandler).Methods("POST")
	r.HandleFunc("/receipts/{id:[0-9a-f-]+}/points", GetReceiptPointsHandler).Methods("GET")

	port := 4000
	addr := fmt.Sprintf(":%d", port)
	http.Handle("/", r)
	fmt.Printf("Server listening on port %d...\n", port)
	http.ListenAndServe(addr, nil)
}

// jsonMiddleware sets the Content-Type header to JSON for all responses.
func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// isTimeBetween2And4 checks if the given time is between 2:00 PM and 4:00 PM.
func isTimeBetween2And4(timeStr string) (output bool, err error) {
	currentTime, err := time.Parse("15:04", timeStr)
	if err != nil {
		return false, err
	}

	startTime := time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC)
	endTime := time.Date(0, 1, 1, 16, 0, 0, 0, time.UTC)

	if currentTime.After(startTime) && currentTime.Before(endTime) {
		return true, nil
	}
	return false, nil
}

// countAlphanumericCharacters counts the number of alphanumeric characters in a string.
func countAlphanumericCharacters(str string) int {
	count := 0
	for _, char := range str {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
			count++
		}
	}
	return count
}

// isValidDateFormat checks if the provided date is in the yyyy-mm-dd format.
func isValidDateFormat(date string) bool {
	// Regular expression to match yyyy-mm-dd format
	datePattern := `^\d{4}-\d{2}-\d{2}$`

	match, err := regexp.MatchString(datePattern, date)
	if err != nil {
		return false
	}
	return match
}

// getPoints calculates points based on various conditions.
func getPoints(retailer string, purchaseDate string, purchaseTime string, total_str string, items []Item) (pointsFinal int, err error) {
	points := countAlphanumericCharacters(retailer)
	total, err := strconv.ParseFloat(total_str, 64)
	if err != nil {
		return 0, err
	}
	// tests

	if retailer == "" || purchaseDate == "" || purchaseTime == "" || total_str == "" || items == nil {
		return 0, err
	}

	if !isValidDateFormat(purchaseDate) {
		return 0, err
	}

	for _, item := range items {
		if item.ShortDescription == "" || item.Price == "" {
			return 0, err
		}
		_, err2 := strconv.ParseFloat(item.Price, 64)
		if err2 != nil {
			return 0, err
		}
	}

	///////////////


	if total == math.Trunc(total) {
		points += 50
	}

	if math.Mod(total, 0.25) == 0 {
		points += 25
	}

	numItems := len(items)
	points += (numItems / 2) * 5

	for _, item := range items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			pricefloatValue, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				return 0, err
			}
			points += int(math.Ceil(pricefloatValue * 0.2))
		}
	}

	day, _ := strconv.Atoi(purchaseDate[8:])
	if day%2 == 1 {
		points += 6
	}

	between2and4, err := isTimeBetween2And4(purchaseTime)
	if err != nil {
		return 0, err
	}

	if between2and4 {
		points += 10
	}

	return points, nil
}

// ProcessHandler processes the receipt and calculates points.
func ProcessHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Retailer      string  `json:"retailer"`
		PurchaseDate  string  `json:"purchaseDate"`
		PurchaseTime  string  `json:"purchaseTime"`
		Total         string `json:"total"`
		Items         []Item  `json:"items"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	points, err := getPoints(input.Retailer, input.PurchaseDate, input.PurchaseTime, input.Total, input.Items)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	savedReceipts[id] = points

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": id,
	})
}

// GetReceiptPointsHandler retrieves the calculated points for a given receipt ID.
func GetReceiptPointsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	points, exists := savedReceipts[id]
	if !exists {
		http.NotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"points": points,
	})
}
