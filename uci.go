package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// parse user/gui move string input e.g. "f7f8q"
func parseMove(moveString string) uint64 {
	var moveList []uint64
	moveList = generateAllMoves(moveList)
	sourceSquare := moveString[:2]
	targetSquare := moveString[2:4]

	// loop through all the possible moves and find the given move
	for _, move := range moveList {
		if sourceSquare == algebraic[getMoveAttr(move, "source")] &&
			targetSquare == algebraic[getMoveAttr(move, "target")] {

			// debug
			//fmt.Printf("debug - source: %d target: %d\n", getMoveAttr(move, "source"), getMoveAttr(move, "target"))
			//fmt.Printf("debug - source: %s target: %s\n", algebraic[getMoveAttr(move, "source")], algebraic[getMoveAttr(move, "target")])

			promotedPiece := getMoveAttr(move, "promoted")

			// check if a piece was promoted
			if promotedPiece != 0 {

				// queen promotion
				if (promotedPiece == Q || promotedPiece == q) && (moveString[4] == 'q') {
					return move
				}

				// rook promotion
				if (promotedPiece == R || promotedPiece == r) && (moveString[4] == 'r') {
					return move
				}

				// bishop promotion
				if (promotedPiece == B || promotedPiece == b) && (moveString[4] == 'b') {
					return move
				}

				// knight promotion
				if (promotedPiece == N || promotedPiece == n) && (moveString[4] == 'n') {
					return move
				}
				continue
			}
			return move
		}
	}
	return 0
}

/*
   Example UCI commands to init position on chess board

   // init start position
   position startpos

   // init start position and make the moves on chess board
   position startpos moves e2e4 e7e5

   // init position from FEN string
   position fen r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1

   // init position from fen string and make moves on chess board
   position fen r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1 moves e2a6 e8g8
*/

// parse UCI "position" command
func parsePosition(command string) bool {
	moveFlag := false

	// split the command by whitespace and loop through words
	words := strings.Fields(command)
	for i, word := range words {
		// command uses either "startpos" or "fen" to tell engine what position to init
		if i == 1 {
			if word == "startpos" {
				startFEN, err := FENfromString(startPosition)
				if err != nil {
					log.Fatal(err)
				} else {
					parseFEN(startFEN)
				}

			} else if word == "fen" {
				FENStringSlice := words[i+1 : i+7]
				FENString := strings.Join(FENStringSlice, " ")
				FENType, err := FENfromString(FENString)
				if err != nil {
					log.Fatal(err)
				} else {
					parseFEN(FENType)
				}
			}
		}

		if moveFlag {
			// check for valid move formatting and if left in check
			if parseMove(word) == 0 || makeMove(parseMove(word)) == 0 {
				return false
			}
		}

		if word == "moves" {
			moveFlag = true
		}

	}
	return true
}

/*
   Example UCI commands to make engine search for the best move

   // fixed depth search
   go depth 64
*/
// parse UCI "go" command
func parseGo(command string) uint64 {
	words := strings.Fields(command)
	depth := words[2]
	intDepth, err := strconv.Atoi(depth)
	if err == nil {
		searchPosition(intDepth)
	} else {
		fmt.Printf("error: %s", err)
	}
	return bestMove
}

func uciLoop() {
	for {
		fmt.Print("\nenter command: ")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		command := strings.Fields(input.Text())[0]

		// parse UCI "isready" command
		if command == "isready" {
			fmt.Printf("readyok\n")
			continue
		} else if command == "position" {
			parsePosition(input.Text())
		} else if command == "ucinewgame" {
			parsePosition("position startpos")
		} else if command == "go" {
			bestMove = parseGo(input.Text())
			printMove(bestMove)

		} else if command == "quit" {
			break
		} else if command == "uci" {
			fmt.Printf("id name goChess\n")
			fmt.Printf("uciok\n")
		}
	}
}
