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

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/", mapHandler)
	router.HandleFunc("/api", apiHandler).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./assets/")))
	log.Fatal(http.ListenAndServe(":8080", router))
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

	sqlStatement := `
			SELECT latitude, longitude, log(frequency) FROM addresses WHERE
			latitude < ($1) AND latitude > ($2)
			AND longitude > ($3) AND longitude < ($4)
			`
	rows, err := db.Query(sqlStatement, north, south, west, east)

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

type Location struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Frequency float32 `json:"frequency"`
}
