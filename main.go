package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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
	keys, ok := r.URL.Query()["FEN"]

	// check that fen string is present
	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'FEN' is missing")
		return
	}

	// log FEN
	FEN := keys[0]
	log.Println("Url Param 'FEN' is: " + FEN)
	log.Println("String replaced FEN is: " + strings.Replace(FEN, "_", " ", -1))
	formattedFEN := "position fen " + strings.Replace(FEN, "_", " ", -1)

	// check legality
	isFENLegal := parsePosition(formattedFEN)

	// build json struct
	response := FENjson{FEN: FEN, legal: isFENLegal}

	// send json response
	json.NewEncoder(w).Encode(response)
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

	debug := false
	if debug {
		FEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR_w_KQkq_-_0_1_moves_e2e4"
		fmt.Println(parsePosition(strings.Replace(FEN, "_", " ", -1)))
	} else {
		// get port
		port := os.Getenv("PORT")

		if port == "" {
			log.Fatal("$PORT must be set")
		}

		// handler
		handleRequests(port)
	}

}
