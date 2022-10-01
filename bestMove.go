package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strings"
)

var DEPTH = 6

func findBestMove(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This endpoint parses the FEN string and returns the best move")
	fmt.Println("Endpoint Hit: bestMove")

	// testing getting GET parameters
	keys, ok := r.URL.Query()["FEN"]

	// check that fen string is present
	if !ok || len(keys[0]) < 1 {
		log.Println("URL parameter 'FEN' is missing")
		return
	}

	// log FEN
	FEN := keys[0]
	log.Println("URL parameter 'FEN' is: " + FEN)
	log.Println("String replaced FEN is: " + strings.Replace(FEN, "_", " ", -1))
	formattedFEN := "position fen " + strings.Replace(FEN, "_", " ", -1)

	// check legality
	isFENLegal := parsePosition(formattedFEN)
	log.Printf("parsePosition(formattedFEN): %t", isFENLegal)

	// find best move
	bestMove = searchPosition(DEPTH)

	// build json struct
	response := FENjson{FEN: FEN}

	// send json response
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error while sending json response %s", err)
	} else {
		log.Printf("JSON encoding succesful - %s", FEN)
	}

}

//type bookMove struct {
//	FEN string
//}

// look at the opening_book db table and see if we can just return the best move from there instead of searching/evaluating manually
//func checkForBookMove(FEN string) string {
//	// Create DB pool
//	DatabaseURL := os.Getenv("DATABASE_URL")
//	db, err := sql.Open("postgres", DatabaseURL)
//	if err != nil {
//		log.Fatal("Failed to open a DB connection: ", err)
//	}
//	defer db.Close()
//
//	// Create an empty best move struct and make the sql query (using $1 for the parameter)
//	var topBookMove bookMove
//	query := "SELECT move_1 FROM opening_book WHERE fen = $1"
//	err = db.QueryRow(query, FEN).Scan(&topBookMove.FEN)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			log.Printf("No book move found for FEN: %s", FEN) // debug
//			return "No book move found"
//		} else {
//			log.Fatal("Failed to execute query: ", err)
//		}
//	}
//
//	// return FEN with the best move made on the board
//	fmt.Printf("Book move result: %s\n", topBookMove.FEN)
//	return topBookMove.FEN
//}
