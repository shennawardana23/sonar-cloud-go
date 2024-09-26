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

func init() {
	var err error
	// Replace with your actual MySQL connection details
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db_pba_local")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/user", getUser)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getUser(w http.ResponseWriter, r *http.Request) {
	first_name := r.URL.Query().Get("first_name")
	last_name := r.URL.Query().Get("last_name")

	// Vulnerable SQL query (potential SQL injection)
	query := fmt.Sprintf("SELECT id, name FROM users WHERE first_name = '%s' AND last_name = '%s'", first_name, last_name)

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
		users = append(users, fmt.Sprintf("User ID: %d, Name: %s", id, name))
	}

	if len(users) == 0 {
		fmt.Fprintf(w, "No users found")
	} else {
		fmt.Fprintf(w, strings.Join(users, "\n"))
	}
}
