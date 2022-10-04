package main

import (
	"errors"
	"strconv"
	"strings"
)

var halfMoves int
var fullMoves int

type FEN struct {
	FENStr       string
	piecesStr    string
	turnStr      string
	castlingStr  string
	enPassantStr string
	halfMovesStr string
	fullMovesStr string
}

func FENfromString(FENstring string) (FEN, error) {
	// FEN strings are formatted with spaces separating information about the position
	// 1st group of characters is the ranks and their piecesStr/spaces
	// 2nd group is the active color
	// 3rd is castling rights
	// 4th is possible en-passant targets
	// 5th is half moves
	// 6th is full moves
	splitFEN := strings.Split(FENstring, " ")
	if len(splitFEN) != 6 {
		return FEN{"", "", "", "", "", "", ""},
			errors.New("malformed FEN string")
	}

	pieces := splitFEN[0]
	turn := splitFEN[1]
	castling := splitFEN[2]
	enPassant := splitFEN[3]
	halfMoves := splitFEN[4]
	fullMoves := splitFEN[5]

	return FEN{FENstring, pieces, turn, castling, enPassant, halfMoves, fullMoves}, nil
}

// generates FEN string of a game state
func generateFEN() (FEN, error) {
	// FEN strings are formatted with spaces separating information about the position
	// 1st group of characters is the ranks and their pieces/spaces
	FENstring := ""

	// set pieces
	pieces := make([]string, 64)
	for piece := P; piece <= k; piece++ { // debug
		bitboard := bitboards[piece]
		for bitboard != 0 {
			index := getLeastSignificantBitIndex(bitboard)
			pieces[index] = stringPieces[piece]
			bitboard = popBit(bitboard, index)
		}
	}

	// set empty squares
	bitboard := ^occupancies[both]
	for bitboard != 0 {
		index := getLeastSignificantBitIndex(bitboard)
		pieces[index] = "-"
		bitboard = popBit(bitboard, index)
	}

	// build FEN string by ranging over pieces/empty squares string and inserting slashes and numbers
	pieceString := strings.Join(pieces, "")
	emptyCount := 0
	for pos, char := range pieceString {
		// handle if there is a piece
		if string(char) != "-" {
			if emptyCount != 0 {
				FENstring += strconv.Itoa(emptyCount)
			}
			emptyCount = 0
			FENstring += string(char)
		} else { // add empty spaces until you reach a piece or end of row
			emptyCount++
		}
		// add slashes at end of rows
		if (pos+1)%8 == 0 && (0 < pos) && (pos < 63) {
			if emptyCount != 0 {
				FENstring += strconv.Itoa(emptyCount)
			}
			emptyCount = 0
			FENstring += "/"
		}
	}

	// 2nd group is the active color
	activeColor := " w "
	if side == black {
		activeColor = " b "
	}
	// 3rd is castling rights
	// 4th is possible en-passant targets
	// 5th is half moves
	// 6th is full moves
	FENstring += activeColor + castleToString(castle) + " " + algebraic[enPassantSquare] + " " + strconv.Itoa(halfMoves) + " " + strconv.Itoa(fullMoves)

	return FENfromString(FENstring)
}

// parse FEN string
func parseFEN(FEN FEN) {
	// reset variables
	bitboards = [12]uint64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	occupancies = [3]uint64{0, 0, 0}
	side = 0
	enPassantSquare = noSquare
	castle = 0
	halfMoves = 0
	fullMoves = 0

	// loop through all the squares and look at the FEN string to determine which piece, if any, goes there
	FENoffset := uint64(0)
	for rank := uint64(0); rank < 8; rank++ {
		for file := uint64(0); file < 8; file++ {
			square := rank*8 + file
			for {
				FENchar := FEN.piecesStr[FENoffset]
				if (FENchar >= 'a' && FENchar <= 'z') || (FENchar >= 'A' && FENchar <= 'Z') {
					// get the piece
					piece := charPieces[FENchar]

					// place the piece on the square
					bitboards[piece] = setBit(bitboards[piece], square)
					break
				} else if FENchar != '/' {
					blankSpaces, _ := strconv.Atoi(string(FENchar))
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
	if FEN.turnStr == "w" || FEN.turnStr == "W" {
		side = white
	} else {
		side = black
	}

	// castling rights
	for _, char := range FEN.castlingStr {
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
	enPassantSquare = squareStringToUint64(FEN.enPassantStr)

	// increment moves
	halfMoves, _ = strconv.Atoi(FEN.halfMovesStr)
	fullMoves, _ = strconv.Atoi(FEN.fullMovesStr)
	if side == black {
		fullMoves++
	}

	// set occupancies
	// loop through white piecesStr
	for piece := P; piece <= K; piece++ {
		occupancies[white] |= bitboards[piece]
	}

	// loop through black piecesStr
	for piece := p; piece <= k; piece++ {
		occupancies[black] |= bitboards[piece]
	}

	// both
	occupancies[both] = occupancies[white] | occupancies[black]
}
