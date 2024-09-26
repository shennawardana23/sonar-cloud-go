package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

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
	username := r.URL.Query().Get("username")

	// Vulnerable SQL query (potential SQL injection)
	query := fmt.Sprintf("SELECT * FROM users WHERE first_name = '%s'", username)

	var id int
	var name string
	err := db.QueryRow(query).Scan(&id, &name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User ID: %d, Name: %s", id, name)
}
