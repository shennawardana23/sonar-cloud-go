package main

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var secretKey string = "mySecretKey" // Security hotspot: Hardcoded secret

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
	r := gin.Default()

	r.GET("/user/code-smell", getUserCodeSmell)
	r.GET("/user/code-smell-duplicate", getUserCodeSmellDuplicate)
	r.GET("/user/bug", getUserBug)
	r.GET("/user/vulnerability", getUserVulnerability)
	r.GET("/user/security-hotspot", getUserSecurityHotspot)

	log.Fatal(r.Run(":8080"))
}

// Function demonstrating code smells
func getUserCodeSmell(c *gin.Context) {
	firstName := c.Query("first_name")

	// Code smell: Complex boolean expression
	if firstName != "" && firstName != "admin" && firstName != "guest" && firstName != "user" {
		// Poor practice: using string concatenation for SQL query
		query := "SELECT * FROM users WHERE first_name = '" + firstName + "'"

		rows, err := db.Query(query)
		if err != nil {
			// Poor error handling: generic error message
			c.String(http.StatusInternalServerError, "Error occurred")
			return
		}
		defer rows.Close()

		var users []string
		for rows.Next() {
			var id int
			var name string
			// Potential issue: not checking for errors in rows.Scan
			rows.Scan(&id, &name)
			users = append(users, fmt.Sprintf(userOutputFormat, id, name))
		}

		// Code smell: not handling the case where users slice is empty
		c.String(http.StatusOK, strings.Join(users, "\n"))
	}
}

// Duplicate function demonstrating similar logic (code duplication)
func getUserCodeSmellDuplicate(c *gin.Context) {
	firstName := c.Query("first_name")

	// Code smell: Complex boolean expression (duplicated)
	if firstName != "" && firstName != "admin" && firstName != "guest" && firstName != "user" {
		// Poor practice: using string concatenation for SQL query
		query := "SELECT * FROM users WHERE first_name = '" + firstName + "'"

		rows, err := db.Query(query)
		if err != nil {
			// Poor error handling: generic error message
			c.String(http.StatusInternalServerError, "Error occurred")
			return
		}
		defer rows.Close()

		var users []string
		for rows.Next() {
			var id int
			var name string
			// Potential issue: not checking for errors in rows.Scan
			rows.Scan(&id, &name)
			users = append(users, fmt.Sprintf(userOutputFormat, id, name))
		}

		// Code smell: not handling the case where users slice is empty
		c.String(http.StatusOK, strings.Join(users, "\n"))
	}
}

// Function demonstrating potential bugs
func getUserBug(c *gin.Context) {
	firstName := c.Query("first_name")

	// Bug: Incorrect string comparison
	if firstName == "admin" || firstName == "Admin" {
		c.String(http.StatusOK, "Admin user")
		return
	}

	// Using a vulnerable query
	query := fmt.Sprintf("SELECT * FROM users WHERE first_name = '%s'", firstName)

	rows, err := db.Query(query)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		users = append(users, fmt.Sprintf(userOutputFormat, id, name))
	}

	// Bug: Incorrect nil check for slice
	if users == nil {
		c.String(http.StatusOK, "No users found")
	} else {
		c.String(http.StatusOK, strings.Join(users, "\n"))
	}
}

// Function demonstrating SQL injection vulnerability
func getUserVulnerability(c *gin.Context) {
	firstName := c.Query("first_name")

	// Vulnerable SQL query without input validation
	query := fmt.Sprintf("SELECT * FROM users WHERE first_name = '%s' OR 1=1", firstName)

	rows, err := db.Query(query)
	if err != nil {
		c.String(http.StatusInternalServerError, "Database query failed")
		return
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			c.String(http.StatusInternalServerError, "Error scanning results")
			return
		}
		users = append(users, fmt.Sprintf(userOutputFormat, id, name))
	}

	c.String(http.StatusOK, strings.Join(users, "\n"))
}

// Function demonstrating security hotspots and other issues
func getUserSecurityHotspot(c *gin.Context) {
	firstName := c.Query("first_name")
	password := c.Query("password")

	// Security hotspot: Weak cryptographic algorithm
	hashedPassword := fmt.Sprintf("%x", md5.Sum([]byte(password)))

	// Potential SQL injection (repeated for emphasis)
	query := fmt.Sprintf("SELECT * FROM users WHERE first_name = '%s' AND password = '%s'", firstName, hashedPassword)

	rows, err := db.Query(query)
	if err != nil {
		// Security issue: Detailed error exposed to user
		c.String(http.StatusInternalServerError, fmt.Sprintf("Database error: %v", err))
		return
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var id int
		var name string
		var storedPassword string
		if err := rows.Scan(&id, &name, &storedPassword); err != nil {
			c.String(http.StatusInternalServerError, "Error scanning results")
			return
		}
		// Security issue: Timing attack vulnerability
		if storedPassword == hashedPassword {
			users = append(users, fmt.Sprintf(userOutputFormat, id, name))
		}
	}

	// Security hotspot: Information exposure
	if len(users) == 0 {
		c.String(http.StatusOK, "Invalid username or password")
	} else {
		c.String(http.StatusOK, strings.Join(users, "\n"))
	}
}
