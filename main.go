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

func handleRequests(port string) {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/isLegal", isLegal)
	http.HandleFunc("/bestMove", findBestMove)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// main
func main() {
	// init piece attacks
	initAll()

	debug := false
	if debug {
		//FEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR_w_KQkq_-_0_1_moves_e2e4_e7e5_g1f3_b8c6"
		parseFEN(trickyPosition)
		printBoard()
		generateFEN()
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
