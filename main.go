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
	log.Printf("Successfully connected to MySQL database at %s:%s", dbHost, dbPort)

	http.HandleFunc("/current-time", currentTimeHandler)
	http.HandleFunc("/request-logs", RequestLogsHandler)
	log.Printf("Server starting on port 80...")
	log.Fatal(http.ListenAndServe(":80", nil))
}

func currentTimeHandler(w http.ResponseWriter, r *http.Request) {
	// Check if method is GET
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		log.Printf("Method %s not allowed from %s", r.Method, r.RemoteAddr)
		return
	}

	log.Printf("Received GET request from %s for current time", r.RemoteAddr)

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
	log.Printf("Stored time %v in database", currentTime)

	// Create the response
	response := TimeResponse{
		CurrentTime: currentTime.Format("2006-01-02 15:04:05 MST") + " (Toronto Time)",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func RequestLogsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		log.Printf("Method %s not allowed from %s", r.Method, r.RemoteAddr)
		return
	}

	log.Printf("Received GET request from %s for request logs", r.RemoteAddr)

	// Check if any records exist
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM time_log").Scan(&count)
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		log.Println("Database query error:", err)
		return
	}

	// If no records found, return appropriate message
	if count == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "No time requests recorded yet",
		})
		return
	}

	// Get Toronto location
	loc, err := time.LoadLocation("America/Toronto")
	if err != nil {
		http.Error(w, "Time zone conversion error", http.StatusInternalServerError)
		log.Println("Time zone conversion error:", err)
		return
	}

	// Query database with correct timezone offset for Toronto
	query := `
        SELECT DATE_FORMAT(
            timestamp,
            '%Y-%m-%d %H:%i:%s'
        ) 
        FROM time_log 
        ORDER BY timestamp DESC
    `
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		log.Println("Database query error:", err)
		return
	}
	defer rows.Close()

	// Store timestamps in an array
	var timestamps []string
	loc, _ = time.LoadLocation("America/Toronto")

	// Iterate through rows and convert to Toronto time
	for rows.Next() {
		var timeStr string
		if err := rows.Scan(&timeStr); err != nil {
			http.Error(w, "Database scan error", http.StatusInternalServerError)
			log.Println("Database scan error:", err)
			return
		}

		// Parse time string and convert to Toronto time
		utcTime, err := time.Parse("2006-01-02 15:04:05", timeStr)
		if err != nil {
			http.Error(w, "Time parsing error", http.StatusInternalServerError)
			log.Println("Time parsing error:", err)
			return
		}

		// Convert UTC to Toronto time
		torontoTime := utcTime.In(loc)
		formattedTime := torontoTime.Format("2006-01-02 15:04:05 MST") + " (Toronto Time)"
		timestamps = append(timestamps, formattedTime)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(timestamps)
}
