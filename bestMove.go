package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func findBestMove(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This endpoint parses the FEN string and returns the best move the engine can find")
	fmt.Println("Endpoint Hit: bestMove")

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
	log.Printf("parsePosition(formattedFEN): %t", isFENLegal)

	// find best move
	bestMove = searchPosition(6)

	// build json struct
	response := bestMovejson{FEN: FEN,
		legal: isFENLegal,
		best:  decodeMove(bestMove)}

	log.Printf(decodeMove(bestMove))

	// send json response
	json.NewEncoder(w).Encode(response)
}
