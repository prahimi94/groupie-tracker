package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestHandleIndex(t *testing.T) {
	// // Change to the root directory of the project
	err := os.Chdir("../../")
	if err != nil {
		t.Fatalf("Error changing directory: %v", err)
	}
	// Test for GET request
	log.Println("Starting TestHandleIndex - GET request")
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		log.Fatalf("Failed to create GET request: %v", err)
	}

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleIndex)

	// Perform the GET request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		log.Printf("TestHandleIndex - GET request: Expected status %v, got %v", http.StatusOK, status)
		t.Errorf("HandleIndex returned wrong status code: got %v want %v", status, http.StatusOK)
	} else {
		log.Printf("TestHandleIndex - GET request: Received status %v, as expected", status)
	}

	// Check if the response contains the correct content (assuming the HTML file is served)
	expected := `<p class="col-md-4 mb-0 text-body-secondary">© 2024 GT by PR & MR</p>` // Add a unique identifier that would appear in the HTML

	if !strings.Contains(rr.Body.String(), expected) {
		log.Printf("TestHandleIndex - GET request: Expected content not found in response body")
		t.Errorf("HandleIndex returned unexpected body: got %v want %v", rr.Body.String(), expected)
	} else {
		log.Printf("TestHandleIndex - GET request: Correct content found in response body")
	}

	// Test for POST request (Should return Method Not Allowed)
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		log.Fatalf("Failed to create POST request: %v", err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		log.Printf("TestHandleIndex - POST request: Expected status %v, got %v", http.StatusMethodNotAllowed, status)
		t.Errorf("HandleIndex returned wrong status code for POST: got %v want %v", status, http.StatusMethodNotAllowed)
	} else {
		log.Printf("TestHandleIndex - POST request: Received status %v, as expected", status)
	}
}

func TestHandleArtists(t *testing.T) {
	// Create a mock server to simulate the API
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return mock JSON data for the /artists endpoint
		if r.URL.Path == "/api/artists" {
			mockData := []ArtistsData{
				{
					Id:      1,
					Image:   "https://groupietrackers.herokuapp.com/api/images/queen.jpeg",
					Name:    "Queen",
					Members: []string{"Freddie Mercury", "Brian May", "Roger Taylor", "John Deacon"},
					// Add other fields as needed
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mockData)
		} else {
			http.NotFound(w, r)
		}
	}))
	defer mockServer.Close()

	// Mock the API URL
	apiUrls["artists"] = mockServer.URL + "/api/artists"

	// Test for GET request
	log.Println("Starting TestHandleArtists - GET request")
	req, err := http.NewRequest("GET", "/artists", nil)
	if err != nil {
		t.Fatalf("Failed to create GET request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleArtists)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("HandleArtists returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check if the response contains the correct content (assuming the HTML file is served)
	expected := `<div class="card" style="width: 100%;">
                  <img src="https://groupietrackers.herokuapp.com/api/images/queen.jpeg" class="card-img-top" alt="Queen" style="max-height: 286px;">
                  <div class="card-body">
                    <h5 class="card-title mb-3">Queen</h5>
                    <a href="artist/1" class="btn btn-outline-info">Show More Info</a>
                  </div>
                </div>`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("HandleArtists returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// Test for POST request (Should return Method Not Allowed)
	req, err = http.NewRequest("POST", "/artists", nil) // Correct route for POST request
	if err != nil {
		t.Fatalf("Failed to create POST request: %v", err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("HandleArtists returned wrong status code for POST: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestHandleErrorPage(t *testing.T) {
	testCases := []struct {
		name         string
		errorType    ErrorPageData
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Test 404 Not Found",
			errorType:    NotFoundError,
			expectedCode: http.StatusNotFound,
			expectedBody: "Page not found",
		},
		{
			name:         "Test 400 Bad Request",
			errorType:    BadRequestError,
			expectedCode: http.StatusBadRequest,
			expectedBody: "Bad request",
		},
		{
			name:         "Test 500 Internal Server Error",
			errorType:    InternalServerError,
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Internal server error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()

			handleErrorPage(w, req, tc.errorType)

			// get result
			res := w.Result()
			defer res.Body.Close()

			// check http status code
			if res.StatusCode != tc.expectedCode {
				t.Errorf("Expected status code %d, but got %d", tc.expectedCode, res.StatusCode)
			}

			// check response content
			body := w.Body.String()
			if !strings.Contains(body, tc.expectedBody) {
				t.Errorf("Expected body to contain '%s', but got '%s'", tc.expectedBody, body)
			}
		})
	}
}

func TestMain(m *testing.M) {
	// Set up global variables
	publicUrl = "frontend/public/"
	apiUrls = map[string]string{
		"base":      "https://groupietrackers.herokuapp.com/api",
		"artists":   "https://groupietrackers.herokuapp.com/api/artists",
		"dates":     "https://groupietrackers.herokuapp.com/api/dates",
		"Dates":     "https://groupietrackers.herokuapp.com/api/Dates",
		"relations": "https://groupietrackers.herokuapp.com/api/relations",
	}

	// Run tests
	os.Exit(m.Run())
}
