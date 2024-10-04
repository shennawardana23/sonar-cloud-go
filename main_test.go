package main

// import (
// 	"net/http/httptest"
// 	"strings"
// 	"testing"
// )

// func TestGetUserVulnerable(t *testing.T) {
// 	tests := []struct {
// 		name           string
// 		firstName      string
// 		expectedOutput string
// 	}{
// 		{
// 			name:           "Valid first name",
// 			firstName:      "John",
// 			expectedOutput: "User ID: 1, Name: John\n", // Adjust based on your database
// 		},
// 		{
// 			name:           "SQL Injection - OR 1=1",
// 			firstName:      "' OR 1=1 --",
// 			expectedOutput: "User ID: 1, Name: John\nUser ID: 2, Name: Jane\n", // Adjust based on your database
// 		},
// 		{
// 			name:           "SQL Injection - UNION SELECT",
// 			firstName:      "' UNION SELECT id, name FROM users --",
// 			expectedOutput: "User ID: 1, Name: John\nUser ID: 2, Name: Jane\n", // Adjust based on your database
// 		},
// 		{
// 			name:           "SQL Injection - Comment",
// 			firstName:      "'; DROP TABLE users; --",
// 			expectedOutput: "No users found", // This should ideally not drop the table
// 		},
// 		{
// 			name:           "Empty input",
// 			firstName:      "",
// 			expectedOutput: "No users found",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			req := httptest.NewRequest("GET", "/user/vulnerable?first_name="+tt.firstName, nil)
// 			w := httptest.NewRecorder()

// 			getUserVulnerable(w, req)

// 			// res := w.Result() // This line is removed
// 			body := w.Body.String()

// 			if !strings.Contains(body, tt.expectedOutput) {
// 				t.Errorf("expected %q in response, got %q", tt.expectedOutput, body)
// 			}
// 			// res := w.Result() // This line is removed
// 		})
// 	}
// }
