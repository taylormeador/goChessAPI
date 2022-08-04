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
		log.Println("Url Param 'FEN' is missing")
		return
	}

	// log FEN
	FEN := keys[0]
	log.Println("Url parameter 'FEN' is: " + FEN)
	log.Println("String replaced FEN is: " + strings.Replace(FEN, "_", " ", -1))
	formattedFEN := "position fen " + strings.Replace(FEN, "_", " ", -1)

	// check legality
	isFENLegal := parsePosition(formattedFEN)
	log.Printf("parsePosition(formattedFEN): %t", isFENLegal)

	// build json struct
	response := FENjson{FEN: FEN,
		legal: isFENLegal}

	// send json response
	json.NewEncoder(w).Encode(response)
}
