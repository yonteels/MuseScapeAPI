package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "admin"
	password = "admin"
	dbname   = "Musescapes"
)

var db *sql.DB

func main() {
	// String connection to database
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Connect to the database
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected!")

	// Define HTTP route and handler
	http.HandleFunc("/album", getAlbum)

	log.Println("Server starting on localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

// HTTP handler that retrieves album data from the database
func getAlbum(w http.ResponseWriter, r *http.Request) {

	// SQL Statement
	sqlStatement := `SELECT "Song Name" FROM user_song WHERE "Player ID" = 1;`

	// query
	rows, err := db.Query(sqlStatement)
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Cursor
	for rows.Next() {
		var songName string
		err := rows.Scan(&songName)
		if err != nil {
			http.Error(w, "Error scanning data", http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, songName)
	}

	// Check for errors after finishing iteration
	if err = rows.Err(); err != nil {
		http.Error(w, "Error iterating rows", http.StatusInternalServerError)
	}
}
