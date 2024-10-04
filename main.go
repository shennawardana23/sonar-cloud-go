package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// Define a constant for the user output format
const userOutputFormat = "User ID: %d, Name: %s"

func init() {
	var err error
	// Replace with your actual MySQL connection details
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db_pba_local")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/user/vulnerable", getUserVulnerable)
	// http.HandleFunc("/user/code-smell", getUserCodeSmell)
	// http.HandleFunc("/user/bug", getUserBug)
	// http.HandleFunc("/user/vulnerability", getUserVulnerability)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Vulnerable function demonstrating SQL Injection
func getUserVulnerable(w http.ResponseWriter, r *http.Request) {
	firstName := r.URL.Query().Get("first_name")

	// Vulnerable SQL query (potential SQL injection)
	query := fmt.Sprintf("SELECT * FROM users WHERE first_name = '%s' OR 1=1", firstName)

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, fmt.Sprintf(userOutputFormat, id, name))
	}

	if len(users) == 0 {
		fmt.Fprintf(w, "No users found")
	} else {
		fmt.Fprintf(w, strings.Join(users, "\n"))
	}
}

// // Function demonstrating a code smell with poor error handling
// func getUserCodeSmell(w http.ResponseWriter, r *http.Request) {
// 	firstName := r.URL.Query().Get("first_name")

// 	// Poor error handling
// 	query := "SELECT * FROM users WHERE first_name = '" + firstName + "'"

// 	rows, err := db.Query(query)
// 	if err != nil {
// 		http.Error(w, "Error occurred", http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()

// 	var users []string
// 	for rows.Next() {
// 		var id int
// 		var name string
// 		if err := rows.Scan(&id, &name); err != nil {
// 			http.Error(w, "Error occurred", http.StatusInternalServerError)
// 			return
// 		}
// 		users = append(users, fmt.Sprintf(userOutputFormat, id, name))
// 	}

// 	if len(users) == 0 {
// 		fmt.Fprintf(w, "No users found")
// 	} else {
// 		fmt.Fprintf(w, strings.Join(users, "\n"))
// 	}
// }

// // Function demonstrating a potential bug in handling empty results
// func getUserBug(w http.ResponseWriter, r *http.Request) {
// 	firstName := r.URL.Query().Get("first_name")

// 	// Using a vulnerable query
// 	query := fmt.Sprintf("SELECT * FROM users WHERE first_name = '%s'", firstName)

// 	rows, err := db.Query(query)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()

// 	var users []string
// 	for rows.Next() {
// 		var id int
// 		var name string
// 		if err := rows.Scan(&id, &name); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		users = append(users, fmt.Sprintf(userOutputFormat, id, name))
// 	}

// 	// Potential bug: not handling the case where users is nil
// 	if users == nil { // {{ edit_14 }}
// 		fmt.Fprintf(w, "No users found")
// 	} else {
// 		fmt.Fprintf(w, strings.Join(users, "\n"))
// 	}
// }

// // Function demonstrating a lack of input validation (vulnerability)
// func getUserVulnerability(w http.ResponseWriter, r *http.Request) {
// 	firstName := r.URL.Query().Get("first_name")

// 	// No input validation
// 	query := fmt.Sprintf("SELECT * FROM users WHERE first_name = '%s'", firstName)

// 	rows, err := db.Query(query)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()

// 	var users []string
// 	for rows.Next() {
// 		var id int
// 		var name string
// 		if err := rows.Scan(&id, &name); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		users = append(users, fmt.Sprintf(userOutputFormat, id, name))
// 	}

// 	if len(users) == 0 {
// 		fmt.Fprintf(w, "No users found")
// 	} else {
// 		fmt.Fprintf(w, strings.Join(users, "\n"))
// 	}
// }
