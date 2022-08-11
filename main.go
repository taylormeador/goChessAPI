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
	fmt.Fprintf(w, "# goChessAPI\n\nThis is a branch from my main goChess project. My intention is to have the engine act as an API that a chess website or GUI can make requests to in order to play chess vs an engine.\n\n\n\nHow to use this API:\n\nThe engine currently has two end points that serve two different functions.\n1. /isLegal - takes a FEN string with a desired move as a GET parameter and returns a JSON response with an updated FEN string. If the desired move is legal, the updated FEN will be the game state after making the move. If the move is illegal, the updated FEN will be the original position. \n\nExample request from client:\ngo-chess-api.herokuapp.com/isLegal?FEN=rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1 moves e2e4\n\nThat is a legal move so the engine would return the FEN in JSON with the move made on the board:\n\"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 1 1\"\n\nIf the client were to incude a move that is deemed not legal by the engine, e.g. \"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1 moves e2e5\",\nthe engine would return \"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1\" which is just the state of the game before the client tryied to make the illegal move.\n\n2. /bestMove - takes a FEN string as a GET paramater and returns a JSON response of the updated FEN with the best move made.\n\nExample request from client:\ngo-chess-api.herokuapp.com/bestMove?FEN=rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1\n\nThe engine would calculate to a predetermined depth what the best move is according to it's scoring algorithm, then make that move, then return the FEN of the game state with that move made. \n\ne.g. the engine might think e2e4 is the best move from the start position so it could return JSON of \"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 1 1\" for our example request\n\n\n\n\nMore info:\n\nDocumentation on FEN and what constitutes a valid FEN string in chess\nhttps://www.chess.com/terms/fen-chess\n")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests(port string) {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/isLegal", isLegal)
	http.HandleFunc("/bestMove", findBestMove)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

var debug = false

// main
func main() {
	// init piece attacks
	initAll()

	if debug { // testing
		// set env variable for db
		// envVariable()
		//FEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR_w_KQkq_-_0_1_moves_e2e4_e7e5_g1f3_b8c6"
		parseFEN(startPosition)
		//printMove(searchPosition(2))
		checkForBookMove(startPosition)

	} else { // production
		// get port
		port := os.Getenv("PORT")

		if port == "" {
			log.Fatal("$PORT must be set")
		}

		// handler
		handleRequests(port)
	}

}
