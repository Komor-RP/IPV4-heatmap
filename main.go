package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	host := os.Getenv("db_host")
	port := os.Getenv("db_port")
	user := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbname := os.Getenv("db_name")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	database, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = database.Ping()
	if err != nil {
		panic(err)
	}
	db = database

}

// http://localhost:8080/api?top=100&bottom=80.2&left=-20&right=200
func main() {

	router := mux.NewRouter()
	router.HandleFunc("/api", apiHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	top := r.URL.Query().Get("top")
	bottom := r.URL.Query().Get("bottom")
	left := r.URL.Query().Get("left")
	right := r.URL.Query().Get("right")
	if len(top) > 0 && len(bottom) > 0 && len(left) > 0 && len(right) > 0 {
		addresses := makeQuery(top, bottom, left, right)
		json.NewEncoder(w).Encode(addresses)
	}
}

func makeQuery(top, bottom, left, right string) []Location {
	var addresses []Location

	sqlStatement := `
			SELECT * FROM addresses WHERE
			longitude < ($1) AND longitude > ($2)
			AND latitude > ($3) AND latitude < ($4)`
	rows, err := db.Query(sqlStatement, top, bottom, left, right)

	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var point Location
		err = rows.Scan(&point.ID, &point.Latitude, &point.Longitude, &point.Frequency)
		if err != nil {
			panic(err)
		}
		addresses = append(addresses, point)
	}

	return addresses
}

type Location struct {
	ID        int64   `json:"id"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Frequency int16   `json:"frequency"`
}
