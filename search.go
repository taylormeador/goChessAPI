package main

// half move counter and bestMove
var ply int
var bestMove uint64

func negaMax(alpha int, beta int, depth int) int {
	// escape condition
	if depth == 0 {
		return evaluate()
	}

	// init vars
	var score int
	var moveList []uint64
	var legalMoves int
	moveList = generateAllMoves(moveList)

	// loop through moves and evaluate
	for count := 0; count < len(moveList); count++ {
		// store state and increment half move counter
		currentState := returnBoardCopy()
		ply++

		// check that move is legal
		if makeMove(moveList[count]) == 0 {
			ply--
			continue
		}

		// increment legal moves counter
		legalMoves++

		// recursive func call
		score = -negaMax(-beta, -alpha, depth-1)

		// restore board and decrement half move counter
		restoreBoardFromCopy(currentState)
		ply--

		// fail-hard beta cutoff
		if score >= beta {
			// node (move) fails high
			return beta
		}

		// update score and associate with best move
		if score > alpha {
			alpha = score

			// update best move iff the half move counter is 0 (root node)
			// debug
			if ply == 0 {
				bestMove = moveList[count]
			}
		}
	}

	// we don't have any legal moves to make in the current postion
	if legalMoves == 0 {
		// see if opponent is in check
		currentKingBitboard := bitboards[k]
		if side == white {
			currentKingBitboard = bitboards[K]
		}
		inCheck := isSquareAttacked(getLeastSignificantBitIndex(currentKingBitboard), side)

		// king is in check
		if inCheck != 0 {
			// return mating score (assuming closest distance to mating position)
			return -49000 + ply
		} else { // king not in check
			// return stalemate score
			return 0
		}
	}
	return score
}

// search position for the best move
func searchPosition(depth int) uint64 {
	negaMax(-50000, 50000, depth)
	return bestMove
}
