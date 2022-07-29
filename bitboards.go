package main

// individual pieces
var bitboards [12]uint64

// by colors
var occupancies [3]uint64

// side to move
var side = -1

// en passant square
var enPassantSquare = noSquare

// castling rights
var castle = -1

// custom type for collecting board info
type gameState struct {
	bitboards       [12]uint64
	occupancies     [3]uint64
	enPassantSquare uint64
	side            int
	castle          int
}

// global for storing game state
var gameStateCopy gameState

// stores a copy of the gamestate to global variable
func copyBoard() {
	var state gameState

	state.bitboards = bitboards
	state.occupancies = occupancies
	state.enPassantSquare = enPassantSquare
	state.side = side
	state.castle = castle

	gameStateCopy = state
}

// returns a copy of the current game state
func returnBoardCopy() gameState {
	var state gameState

	state.bitboards = bitboards
	state.occupancies = occupancies
	state.enPassantSquare = enPassantSquare
	state.side = side
	state.castle = castle

	return state
}

// sets the relevant global vars to reflect the stored gamestate
func restoreBoard() {
	bitboards = gameStateCopy.bitboards
	occupancies = gameStateCopy.occupancies
	enPassantSquare = gameStateCopy.enPassantSquare
	side = gameStateCopy.side
	castle = gameStateCopy.castle
}

// sets the relevant global vars to reflect the gamestate passed as an argument
func restoreBoardFromCopy(copy gameState) {
	bitboards = copy.bitboards
	occupancies = copy.occupancies
	enPassantSquare = copy.enPassantSquare
	side = copy.side
	castle = copy.castle
}
