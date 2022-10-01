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
	makeMove(bestMove)
	FEN = generateFEN()

	// send json response
	response := FENjson{FEN: FEN}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error while sending json response %s", err)
	} else {
		log.Printf("JSON encoding succesful - %s", FEN)
	}

}
