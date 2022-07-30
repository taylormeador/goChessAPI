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
	fmt.Fprintf(w, "Welcome to the home page")
	fmt.Println("Endpoint Hit: homePage")
}

func isLegal(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This endpoint parses the FEN string and returns the updated game state")
	fmt.Println("Endpoint Hit: isLegal")

	// testing getting GET parameters
	keys, ok := r.URL.Query()["key"]

	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'key' is missing")
		return
	}

	// Query()["key"] will return an array of items,
	// we only want the single item.
	key := keys[0]

	log.Println("Url Param 'key' is: " + string(key))

	//json.NewEncoder(w).Encode(FENstring)
}

func handleRequests(port string) {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/isLegal", isLegal)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// main
func main() {
	// init piece attacks
	initAll()

	// get port
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// handler
	handleRequests(port)
}
