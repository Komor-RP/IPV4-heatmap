package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var addresses []Location

// http://localhost:8080/?top=23.5&bottom=50.2&left=20&right=-20
func main() {
	addresses = append(addresses, Location{ID: 1, Latitude: 50.2, Longitude: 20})
	addresses = append(addresses, Location{ID: 2, Latitude: 29, Longitude: -320})
	addresses = append(addresses, Location{ID: 3, Latitude: -45, Longitude: 120})

	router := mux.NewRouter()
	router.HandleFunc("/", testHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))

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
