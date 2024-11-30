package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type TimeResponse struct {
	CurrentTime string `json:"current_time"`
}

var db *sql.DB

func main() {
	var err error
	// Get environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Connect to MySQL database
	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Verify connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	http.HandleFunc("/current-time", currentTimeHandler)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func currentTimeHandler(w http.ResponseWriter, r *http.Request) {
	// Get the current time in Toronto
	loc, err := time.LoadLocation("America/Toronto")
	if err != nil {
		http.Error(w, "Time zone conversion error", http.StatusInternalServerError)
		log.Println("Time zone conversion error:", err)
		return
	}
	currentTime := time.Now().In(loc)

	// Insert the current time into the database
	_, err = db.Exec("INSERT INTO time_log (timestamp) VALUES (?)", currentTime)
	if err != nil {
		http.Error(w, "Database insert error", http.StatusInternalServerError)
		log.Println("Database insert error:", err)
		return
	}

	// Create the response
	response := TimeResponse{
		CurrentTime: currentTime.Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
