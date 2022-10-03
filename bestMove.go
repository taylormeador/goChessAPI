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

	// check legality
	formattedFEN := "position fen " + strings.Replace(FEN, "_", " ", -1)
	isFENLegal := parsePosition(formattedFEN)
	log.Printf("parsePosition(formattedFEN): %t", isFENLegal)
	if isFENLegal == false {
		log.Printf("Illegal FEN - TODO need to figure out what to do in this case")
	}

	// find best move
	bestMove = searchPosition(DEPTH)
	makeMove(bestMove)
	newFEN := generateFEN()

	// check for checkmate/stalemate
	checkmate := false
	stalemate := false
	if newFEN == FEN {
		currentKingBitboard := bitboards[k]
		if side == white {
			currentKingBitboard = bitboards[K]
		}
		// check if king is in check
		inCheck := isSquareAttacked(getLeastSignificantBitIndex(currentKingBitboard), side)
		if inCheck != 0 {
			checkmate = true
		} else { // king not in check
			stalemate = true
		}
	}

	// send json response
	response := FENjson{
		FEN:       newFEN,
		Checkmate: checkmate,
		Stalemate: stalemate,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error while sending json response %s", err)
	} else {
		log.Printf("JSON encoding succesful - %s", FEN)
	}

}
