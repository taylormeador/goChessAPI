package main

import (
	"fmt"
	"time"
)

var nodes uint64

func perftDriver(depth int) {
	// recursion escape condition
	if depth == 0 {
		nodes++
		return
	}

	// generate moves
	var moveList []uint64
	moveList = generateAllMoves(moveList)

	// store original position
	var state gameState
	state.bitboards = bitboards
	state.occupancies = occupancies
	state.enPassantSquare = enPassantSquare
	state.side = side
	state.castle = castle

	for _, move := range moveList {
		copyBoard()
		if makeMove(move) == 0 {
			continue
		}

		// debug
		//printBoard()
		//
		//// get promotion string
		//promoted := ""
		//if getMoveAttr(move, "promoted") != 0 {
		//	promoted = promotedPieces[getMoveAttr(move, "promoted")]
		//}
		//
		//// print move info
		//fmt.Printf("\n    Move: %s%s%s", algebraic[getMoveAttr(move, "source")],
		//	algebraic[getMoveAttr(move, "target")],
		//	promoted)
		//input := bufio.NewScanner(os.Stdin)
		//input.Scan()
		// end debug

		// recursive call
		perftDriver(depth - 1)

		// reset
		gameStateCopy = state
		restoreBoard()
	}
}

func perftTest(depth int) {
	fmt.Printf("\n    Performance test\n")
	var moveList []uint64
	moveList = generateAllMoves(moveList)
	startTime := time.Now()

	// store original position
	var state gameState
	state.bitboards = bitboards
	state.occupancies = occupancies
	state.enPassantSquare = enPassantSquare
	state.side = side
	state.castle = castle

	// loop through moves and search, counting nodes
	for _, move := range moveList {
		if makeMove(move) == 0 {
			continue
		}

		// debug
		//fmt.Printf("\n")
		//fmt.Printf("********** New top level branch ** Move: %s%s **********************",
		//	algebraic[getMoveAttr(move, "source")],
		//	algebraic[getMoveAttr(move, "target")])
		//printBoard()
		// end debug

		cumulativeNodes := nodes

		// recursive func call
		perftDriver(depth - 1)

		oldNodes := nodes - cumulativeNodes

		// get promotion string
		promoted := ""
		if getMoveAttr(move, "promoted") != 0 {
			promoted = promotedPieces[getMoveAttr(move, "promoted")]
		}

		// print move info
		fmt.Printf("\n    Move: %s%s%s    Nodes: %d", algebraic[getMoveAttr(move, "source")],
			algebraic[getMoveAttr(move, "target")],
			promoted,
			oldNodes)

		// reset
		gameStateCopy = state
		restoreBoard()

	}
	// calculate time taken and print results
	endTime := time.Now()
	fmt.Printf("\n\n    Depth: %d", depth)
	fmt.Printf("\n    Nodes: %d", nodes)
	fmt.Printf("\n    Time: %s\n\n", endTime.Sub(startTime).String())
}
