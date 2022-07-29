package main

import (
	"fmt"
	"strconv"
	"strings"
)

// prints uint64 as bitboard
func printBitboard(bitboard uint64) {
	fmt.Println()
	// loop through ranks and files
	for rank := uint64(0); rank < 8; rank++ {
		fmt.Printf(" %d   ", 8-rank)
		for file := uint64(0); file < 8; file++ {
			// convert rank and file to index
			index := rank*8 + file

			// check whether the bit should be on or off
			printBit := 0
			if bitboard&(uint64(1)<<index) != 0 {
				printBit = 1
			}

			// print 1 or 0 based on previously generated bool
			fmt.Printf(" %d ", printBit)
		}
		fmt.Println()
	}
	fmt.Println()
	// print files and bitboard integer value
	fmt.Println("      a  b  c  d  e  f  g  h")
	fmt.Printf("  Bitboard: %d\n\n", bitboard)
}

// print board
func printBoard() {
	fmt.Println()
	for rank := uint64(0); rank < 8; rank++ {
		fmt.Printf(" %d   ", 8-rank)
		for file := uint64(0); file < 8; file++ {
			square := rank*8 + file
			piece := -1
			for bitboardPiece := P; bitboardPiece <= k; bitboardPiece++ {
				if getBit(bitboards[bitboardPiece], square) != 0 {
					piece = bitboardPiece
				}
			}

			if piece == -1 {
				fmt.Printf(" %s ", ".")
			} else {
				fmt.Printf(" %s ", unicodePieces[piece])
			}

		}
		fmt.Println()
	}
	fmt.Println()
	// print files and bitboard integer value
	fmt.Printf("      a  b  c  d  e  f  g  h\n\n")
	fmt.Printf("side: %s\n", colorFromInt(side))
	fmt.Printf("castling: %s\n", castleToString(castle))
	fmt.Printf("en-passant: %s\n", algebraic[enPassantSquare])
}

// parse FEN string
func parseFEN(FEN string) {
	// reset variables
	bitboards = [12]uint64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	occupancies = [3]uint64{0, 0, 0}
	side = 0
	enPassantSquare = noSquare
	castle = 0

	// FEN strings are formatted with spaces separating information about the position\
	// 1st group of characters is the ranks and their pieces/spaces
	// 2nd group is the active color
	// 3rd is castling rights
	// 4th is possible en-passant targets
	// 5th is half moves
	// 6th is full moves
	splitFEN := strings.Split(FEN, " ")
	if len(splitFEN) != 6 {
		fmt.Println("********** Malformed FEN string ************")
		return
	}

	pieces := splitFEN[0]
	turn := splitFEN[1]
	castling := splitFEN[2]
	enPassant := splitFEN[3]
	//halfMove := splitFEN[4]
	//fullMove := splitFEN[5]

	// loop through all the squares and look at the FEN string to determine which piece, if any, goes there
	FENoffset := uint64(0)
	for rank := uint64(0); rank < 8; rank++ {
		for file := uint64(0); file < 8; file++ {
			square := rank*8 + file
			for {
				FEN := pieces[FENoffset]
				if (FEN >= 'a' && FEN <= 'z') || (FEN >= 'A' && FEN <= 'Z') {
					// get the piece
					piece := charPieces[FEN]

					// place the piece on the square
					bitboards[piece] = setBit(bitboards[piece], square)
					break
				} else if FEN != '/' {
					blankSpaces, _ := strconv.Atoi(string(FEN))
					file += uint64(blankSpaces) - 1
					break
				} else {
					FENoffset += 1
				}
			}
			FENoffset += 1
		}
	}

	// set other board position parameters
	// side to move
	if turn == "w" || turn == "W" {
		side = white
	} else {
		side = black
	}

	// castling rights
	for _, char := range castling {
		switch char {
		case 'K':
			castle += wk
		case 'Q':
			castle += wq
		case 'k':
			castle += bk
		case 'q':
			castle += bq
		}
	}

	// en passant
	enPassantSquare = squareStringToUint64(enPassant)

	// set occupancies
	// loop through white pieces
	for piece := P; piece <= K; piece++ {
		occupancies[white] |= bitboards[piece]
	}

	// loop through black pieces
	for piece := p; piece <= k; piece++ {
		occupancies[black] |= bitboards[piece]
	}

	// both
	occupancies[both] = occupancies[white] | occupancies[black]
}

// print 1 or 0 for a square, depending on if it is attacked by an enemy piece or not
func printAttackedSquares(side int) {
	fmt.Println()
	for rank := uint64(0); rank < 8; rank++ {
		fmt.Printf("  %d  ", 8-rank)
		for file := uint64(0); file < 8; file++ {
			square := rank*8 + file
			fmt.Printf(" %d ", isSquareAttacked(square, side))
		}
		fmt.Println()
	}
	fmt.Println()
	fmt.Print("      a  b  c  d  e  f  g  h\n\n")
}
