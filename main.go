package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Mhmm")
}

var db *sql.DB

func init() {
	var dbErr error
	db, dbErr = sql.Open("sqlite3", "./data.db")
	if dbErr != nil {
		panic(dbErr)
	}

	//Create table if not exists
	_, dbErr = db.Exec("CREATE TABLE IF NOT EXISTS sensors(ID INTEGER PRIMARY KEY AUTOINCREMENT, AverageValue FLOAT, TotalNumberOfEntries INTEGER)")
	if dbErr != nil {
		panic(dbErr)
	}
}

func receiveSensorPayload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	fmt.Println(string(body))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Json Received"))
}

func sendSensorStatus(w http.ResponseWriter, r *http.Request) {

}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/receiveSensorPayload", receiveSensorPayload)
	fmt.Println("Listening on port 8080..")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server", err)
	}

}
