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

type Item struct {
	ShortDescription string  `json:"shortDescription"`
	Price            string `json:"price"`
}

var savedReceipts = make(map[string]int)

func main() {
	r := mux.NewRouter()
	r.Use(jsonMiddleware)
	r.HandleFunc("/receipts/process", ProcessHandler).Methods("POST")
	r.HandleFunc("/receipts/{id:[0-9a-f-]+}/points", GetReceiptPointsHandler).Methods("GET")

	port := 4000
	addr := fmt.Sprintf(":%d", port)
	http.Handle("/", r)
	fmt.Printf("Server listening on port %d...\n", port)
	http.ListenAndServe(addr, nil)
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func isTimeBetween2And4(timeStr string) (output bool, err error) {
	currentTime, err := time.Parse("15:04", timeStr)
	if err != nil {
		return false, err
	}

	startTime := time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC)
	endTime := time.Date(0, 1, 1, 16, 0, 0, 0, time.UTC)

	if (currentTime.After(startTime) && currentTime.Before(endTime)) {
		return true, nil
	}
	return false, nil
}

func countAlphanumericCharacters(str string) int {
	count := 0
	for _, char := range str {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
			count++
		}
	}
	return count
}

func isValidDateFormat(date string) bool {
	// Regular expression to match yyyy-mm-dd format
	datePattern := `^\d{4}-\d{2}-\d{2}$`

	match, err := regexp.MatchString(datePattern, date)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	return match
}

func getPoints(retailer string, purchaseDate string, purchaseTime string, total float64, items []Item) (pointsFinal int, err error) {
	points := countAlphanumericCharacters(retailer)
	fmt.Println("added points for retailer: ", points)
	if total == math.Trunc(total)  {
		fmt.Println("adding 50 points for whole number ")
		points += 50
	}

	if math.Mod(total, 0.25) == 0 {
		fmt.Println("adding 25 points for quarter ")
		points += 25
	}

	numItems := len(items)
	points += (numItems / 2) * 5
	fmt.Println("adding 5 points for every 2 items: ", (numItems / 2) * 5)

	for _, item := range items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			pricefloatValue, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				fmt.Println("Error:", err)
				return 0, err
			}
			points += int(math.Ceil(pricefloatValue * 0.2))
			fmt.Println("adding 20 p of price for item: ", item.ShortDescription, " which is ", int(pricefloatValue * 0.2))
		}
	}

	day, _ := strconv.Atoi(purchaseDate[8:])
	if day%2 == 1 {
		fmt.Println("adding 6 points for odd day ")
		points += 6
	}

	between2and4, err := isTimeBetween2And4(purchaseTime)
	fmt.Println("isTimeBetween2And4: ", between2and4)

	if err != nil {
		fmt.Println("Error: Invalid time format")
		return 0, err
	}

	if between2and4 {
		fmt.Println("adding 10 points for time between 2 and 4 ")
		points += 10
	}

	return points, nil
}

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

	fmt.Println("Input: ", input)

	if (input.Retailer == "" || input.PurchaseDate == "" || input.PurchaseTime == "" || input.Total == "" || input.Items == nil) {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if (!isValidDateFormat(input.PurchaseDate)) {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	for _, item := range input.Items {
		if (item.ShortDescription == "" || item.Price == "") {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
	}


	floatValue, err := strconv.ParseFloat( input.Total, 64)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	points, err := getPoints(input.Retailer, input.PurchaseDate, input.PurchaseTime, floatValue, input.Items)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	savedReceipts[id] = points

	fmt.Println("Points:", points)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": id,
	})
}

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
