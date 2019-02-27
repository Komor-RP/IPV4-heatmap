package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var templates = template.Must(template.ParseFiles("assets/index.html"))
var db *sql.DB

func init() {
	var host, port, user, password, dbname string
	var database *sql.DB
	var err error

	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()

		if err != nil {
			panic(err)
		}

		host = os.Getenv("db_host")
		port = os.Getenv("db_port")
		user = os.Getenv("db_user")
		password = os.Getenv("db_pass")
		dbname = os.Getenv("db_name")

		psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)

		database, err = sql.Open("postgres", psqlInfo)
		if err != nil {
			panic(err)
		}
	} else {
		database, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
		if err != nil {
			panic(err)
		}
	}

	err = database.Ping()
	if err != nil {
		panic(err)
	}
	db = database

}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := mux.NewRouter()
	router.HandleFunc("/", mapHandler)
	router.HandleFunc("/api", apiHandler).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./assets/")))
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func mapHandler(w http.ResponseWriter, r *http.Request) {

	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	north := r.URL.Query().Get("north")
	south := r.URL.Query().Get("south")
	west := r.URL.Query().Get("west")
	east := r.URL.Query().Get("east")
	if len(north) > 0 && len(south) > 0 && len(west) > 0 && len(east) > 0 {
		geojson := makeQuery(north, south, west, east)
		json.NewEncoder(w).Encode(geojson)
	}
}

func makeQuery(north, south, west, east string) []Location {
	var addresses []Location

	maxFreq := getMaxFrequency(north, south, west, east)

	sqlStatement := `
			SELECT latitude, longitude, log(frequency) / log($5)
			FROM addresses WHERE
			latitude < ($1) AND latitude > ($2)
			AND longitude > ($3) AND longitude < ($4)
			`
	rows, err := db.Query(sqlStatement, north, south, west, east, maxFreq)

	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var point Location

		err = rows.Scan(&point.Latitude, &point.Longitude, &point.Frequency)
		if err != nil {
			panic(err)
		}

		addresses = append(addresses, point)
	}

	return addresses
}

func getMaxFrequency(north, south, west, east string) int64 {
	var max int64

	sqlStatement := `SELECT MAX(frequency) from Addresses
		WHERE
		latitude < ($1) AND latitude > ($2)
		AND longitude > ($3) AND longitude < ($4)`
	answer, err := db.Query(sqlStatement, north, south, west, east)

	if err != nil {
		panic(err)
	}

	answer.Next()
	err = answer.Scan(&max)
	if err != nil {
		panic(err)
	}

	return max
}

type Location struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Frequency float32 `json:"frequency"`
}
