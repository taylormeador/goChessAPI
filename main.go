package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// init
func initAll() {
	initLeapersAttacks()
	initSlidersAttacks(bishop)
	initSlidersAttacks(rook)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests(port string) {
	http.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// main
func main() {
	initAll()
	parseFEN(startPosition)

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	handleRequests(port)
}
