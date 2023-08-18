package main

import (
	"testing"
)

func TestIsTimeBetween2And4(t *testing.T) {
	// Test case: time is between 2:00 PM and 4:00 PM
	result, err := isTimeBetween2And4("15:30")
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if !result {
		t.Errorf("Expected time to be between 2:00 PM and 4:00 PM")
	}

	// Test case: time is not between 2:00 PM and 4:00 PM
	result, err = isTimeBetween2And4("10:00")
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if result {
		t.Errorf("Expected time not to be between 2:00 PM and 4:00 PM")
	}

	// Test case: invalid time format
	_, err = isTimeBetween2And4("15:30:00")
	if err == nil {
		t.Errorf("Expected an error for invalid time format")
	}

    // Test case: invalid time format
	_, err = isTimeBetween2And4("15s:30")
	if err == nil {
		t.Errorf("Expected an error for invalid time format")
	}

    // Test case: invalid time format
    _, err = isTimeBetween2And4("")
	if err == nil {
		t.Errorf("Expected an error for invalid time format")
	}
}

func TestCountAlphanumericCharacters(t *testing.T) {
	count := countAlphanumericCharacters("Hello123")
	if count != 8 {
		t.Errorf("Expected count to be 8, but got %d", count)
	}

    count = countAlphanumericCharacters("Hello&123")
	if count != 8 {
		t.Errorf("Expected count to be 8, but got %d", count)
	}

	count = countAlphanumericCharacters("Spaces are ignored 456")
	if count != 19 {
		t.Errorf("Expected count to be 19, but got %d", count)
	}
}

func TestIsValidDateFormat(t *testing.T) {
	// Test case: valid date format
	valid, err := isValidDateFormat("2023-08-17")
    if err != nil {
        t.Errorf("Expected no error, but got: %v", err)
    }
	if !valid {
		t.Errorf("Expected date format to be valid")
	}

	// Test case: invalid date format
	invalid, err2 := isValidDateFormat("08-17-2023")
    if err2 != nil {
        t.Errorf("Expected no error, but got: %v", err2)
    }
	if invalid {
		t.Errorf("Expected date format to be invalid")
	}
}

func TestGetPoints(t *testing.T) {
	input := struct {
		Retailer      string  `json:"retailer"`
		PurchaseDate  string  `json:"purchaseDate"`
		PurchaseTime  string  `json:"purchaseTime"`
		Total         string `json:"total"`
		Items         []Item  `json:"items"`
	}{
		Retailer:     "M&M Corner Market",
		PurchaseDate: "2022-03-20",
		PurchaseTime: "14:33",
		Total:        "9.00",
		Items: []Item{
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
		},
	}

	expectedPoints := 109

	points, err := getPoints(input.Retailer, input.PurchaseDate, input.PurchaseTime, input.Total, input.Items)
	if err != nil {
		t.Errorf("Expected error to be nil, but got: %v", err)
	}
	if points != expectedPoints {
		t.Errorf("Expected points to be %d, but got %d", expectedPoints, points)
	}

    // Total is in wrong format
	input = struct {
		Retailer      string  `json:"retailer"`
		PurchaseDate  string  `json:"purchaseDate"`
		PurchaseTime  string  `json:"purchaseTime"`
		Total         string `json:"total"`
		Items         []Item  `json:"items"`
	}{
		Retailer:     "M&M Corner Market",
		PurchaseDate: "2022-03-20",
		PurchaseTime: "14:33",
		Total:        "9ws.00",
		Items: []Item{
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
		},
	}

    
	_, err = getPoints(input.Retailer, input.PurchaseDate, input.PurchaseTime, input.Total, input.Items)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}

    // Purchase date is in wrong format
    input = struct {
		Retailer      string  `json:"retailer"`
		PurchaseDate  string  `json:"purchaseDate"`
		PurchaseTime  string  `json:"purchaseTime"`
		Total         string `json:"total"`
		Items         []Item  `json:"items"`
	}{
		Retailer:     "M&M Corner Market",
		PurchaseDate: "2022-03-2074",
		PurchaseTime: "14:33",
		Total:        "9s.00",
		Items: []Item{
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
		},
	}


	_, err = getPoints(input.Retailer, input.PurchaseDate, input.PurchaseTime, input.Total, input.Items)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}


   // Purchase date is in wrong format
    input = struct {
		Retailer      string  `json:"retailer"`
		PurchaseDate  string  `json:"purchaseDate"`
		PurchaseTime  string  `json:"purchaseTime"`
		Total         string `json:"total"`
		Items         []Item  `json:"items"`
	}{
		Retailer:     "M&M Corner Market",
		PurchaseDate: "2022-03-2074",
		PurchaseTime: "14:33",
		Total:        "9.00",
		Items: []Item{
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
		},
	}


	_, err = getPoints(input.Retailer, input.PurchaseDate, input.PurchaseTime, input.Total, input.Items)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
    }

    // Item price is in wrong format
    input = struct {
		Retailer      string  `json:"retailer"`
		PurchaseDate  string  `json:"purchaseDate"`
		PurchaseTime  string  `json:"purchaseTime"`
		Total         string `json:"total"`
		Items         []Item  `json:"items"`
	}{
		Retailer:     "M&M Corner Market",
		PurchaseDate: "2022-03-20",
		PurchaseTime: "14:33",
		Total:        "9.00",
		Items: []Item{
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
			{
				ShortDescription: "Gatorade",
				Price:            "no-price",
			},
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
		},
	}


	_, err = getPoints(input.Retailer, input.PurchaseDate, input.PurchaseTime, input.Total, input.Items)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
	
}
