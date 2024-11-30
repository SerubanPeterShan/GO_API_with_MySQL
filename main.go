package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type TimeResponse struct {
	CurrentTime string `json:"current_time"`
}

func main() {
	http.HandleFunc("/current-time", currentTimeHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func currentTimeHandler(w http.ResponseWriter, r *http.Request) {
	// Connect to the database
	db, err := sql.Open("mysql", "root:my-secret-pw@tcp(54.219.52.229:3306)/my_database")
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		log.Println("Database connection error:", err)
		return
	}
	defer db.Close()

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
