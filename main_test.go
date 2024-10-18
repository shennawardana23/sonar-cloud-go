package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
)

func TestGetUserCodeSmell(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a new mock database
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	r := gin.Default()
	r.GET("/user/code-smell", getUserCodeSmell(mockDB))

	tests := []struct {
		name           string
		firstName      string
		mockQuery      string
		mockRows       *sqlmock.Rows
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid first name",
			firstName:      "John",
			mockQuery:      "SELECT \\* FROM users WHERE first_name = 'John'",
			mockRows:       sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "John"),
			expectedStatus: http.StatusOK,
			expectedBody:   "User ID: 1, Name: John",
		},
		{
			name:           "Empty first name",
			firstName:      "",
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
		{
			name:           "Admin first name",
			firstName:      "admin",
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
		{
			name:           "Non-existent user",
			firstName:      "NonExistent",
			mockQuery:      "SELECT \\* FROM users WHERE first_name = 'NonExistent'",
			mockRows:       sqlmock.NewRows([]string{"id", "name"}),
			expectedStatus: http.StatusOK,
			expectedBody:   "No users found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockQuery != "" {
				mock.ExpectQuery(tt.mockQuery).WillReturnRows(tt.mockRows)
			}

			req, _ := http.NewRequest("GET", "/user/code-smell?first_name="+tt.firstName, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d; got %d", tt.expectedStatus, w.Code)
			}
			if !strings.Contains(w.Body.String(), tt.expectedBody) {
				t.Errorf("expected body to contain %q; got %q", tt.expectedBody, w.Body.String())
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetUserBug(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/user/bug", getUserBug)

	tests := []struct {
		name           string
		firstName      string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Admin user",
			firstName:      "admin",
			expectedStatus: http.StatusOK,
			expectedBody:   "Admin user",
		},
		{
			name:           "Admin user (capitalized)",
			firstName:      "Admin",
			expectedStatus: http.StatusOK,
			expectedBody:   "Admin user",
		},
		{
			name:           "Regular user",
			firstName:      "John",
			expectedStatus: http.StatusOK,
			expectedBody:   "No users found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/user/bug?first_name="+tt.firstName, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d; got %d", tt.expectedStatus, w.Code)
			}
			if !strings.Contains(w.Body.String(), tt.expectedBody) {
				t.Errorf("expected body to contain %q; got %q", tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestGetUserVulnerability(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/user/vulnerability", getUserVulnerability)

	tests := []struct {
		name           string
		firstName      string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid first name",
			firstName:      "John",
			expectedStatus: http.StatusOK,
			expectedBody:   "User ID: 1, Name: John",
		},
		{
			name:           "SQL Injection attempt",
			firstName:      "' OR '1'='1",
			expectedStatus: http.StatusOK,
			expectedBody:   "User ID: 1, Name: John\nUser ID: 2, Name: Jane",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/user/vulnerability?first_name="+tt.firstName, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d; got %d", tt.expectedStatus, w.Code)
			}
			if !strings.Contains(w.Body.String(), tt.expectedBody) {
				t.Errorf("expected body to contain %q; got %q", tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestGetUserSecurityHotspot(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/user/security-hotspot", getUserSecurityHotspot)

	tests := []struct {
		name           string
		firstName      string
		password       string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid credentials",
			firstName:      "John",
			password:       "password123",
			expectedStatus: http.StatusOK,
			expectedBody:   "User ID: 1, Name: John",
		},
		{
			name:           "Invalid credentials",
			firstName:      "John",
			password:       "wrongpassword",
			expectedStatus: http.StatusOK,
			expectedBody:   "Invalid username or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/user/security-hotspot?first_name="+tt.firstName+"&password="+tt.password, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d; got %d", tt.expectedStatus, w.Code)
			}
			if !strings.Contains(w.Body.String(), tt.expectedBody) {
				t.Errorf("expected body to contain %q; got %q", tt.expectedBody, w.Body.String())
			}
		})
	}
}
