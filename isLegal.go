package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func isLegal(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: isLegal")

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
	log.Println("String-replaced FEN is: " + strings.Replace(FEN, "_", " ", -1))
	formattedFEN := "position fen " + strings.Replace(FEN, "_", " ", -1)

	// check legality
	isFENLegal := parsePosition(formattedFEN)
	log.Printf("parsePosition(formattedFEN): %t", isFENLegal)

	// build json struct
	var response FENjson
	newFEN := ""
	if isFENLegal == false { // if move is not legal, return the FEN of the board position without the move made
		FENfields := strings.Fields(FEN)
		newFEN = strings.Join(FENfields[0:6], " ")
		log.Printf("Illegal move. Returning FEN: %s", newFEN)
		response = FENjson{FEN: newFEN}
	} else { // if the move is legal, return the updated FEN of the board with the move made
		newFEN = generateFEN()
		log.Printf("Legal move. Returning FEN: %s", newFEN)
		response = FENjson{FEN: newFEN}
	}

	// send json response
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error while sending json response %s", err)
	} else {
		log.Printf("JSON encoding succesful - %s", newFEN)
	}
}
