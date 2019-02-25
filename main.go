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

var addresses []Location

// http://localhost:8080/?top=23.5&bottom=50.2&left=20&right=-20
func main() {
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

	fmt.Println(psqlInfo)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	sqlStatement := `
	SELECT * FROM addresses WHERE
	latitude > ($1) AND latitude < ($2)`
	rows, err := db.Query(sqlStatement, 20, 30)

	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var point Location
		err = rows.Scan(&point.ID, &point.Latitude, &point.Longitude)
		if err != nil {
			panic(err)
		}
		addresses = append(addresses, point)
	}

	fmt.Println(addresses[0])

	router := mux.NewRouter()
	router.HandleFunc("/", testHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))

}

func makeQuery(top, bottom, left, right string) {
	// sqlStatement := `
	// SELECT * FROM addresses WHERE
	// latitude > ($1)`
	// _, err = db.Exec(sqlStatement, 20)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Query())
	w.Header().Set("Content-Type", "application/json")

	top := r.URL.Query().Get("top")
	bottom := r.URL.Query().Get("bottom")
	left := r.URL.Query().Get("left")
	right := r.URL.Query().Get("right")
	if len(top) > 0 && len(bottom) > 0 && len(left) > 0 && len(right) > 0 {
		json.NewEncoder(w).Encode(addresses)
	}
}

type Location struct {
	ID        int64   `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
