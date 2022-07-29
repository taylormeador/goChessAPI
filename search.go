package main

// half move counter
var ply int

// init
var bestMoveSoFar uint64
var bestMove uint64

// negamax with alpha beta pruning
//func negaMax(alpha int, beta int, depth int) int {
//	// recursion escape condition
//	if depth == 0 {
//		return evaluate()
//	}
//
//	// increment node counter
//	nodes++
//
//	// init check vars
//	enemySide := white
//	currentKingBitboard := bitboards[K]
//	if side == white {
//		enemySide = black
//		currentKingBitboard = bitboards[k]
//	}
//
//	// is king in check
//	inCheck := isSquareAttacked(getLeastSignificantBitIndex(currentKingBitboard), enemySide)
//
//	// init vars
//	legalMoves := 0
//	bestMoveSoFar = 0
//	oldAlpha := alpha
//	generateAllMoves()
//
//	// loop over moves within a move list
//	for count := 0; count < len(moveList); count++ {
//
//		// preserve board state
//		currentState := returnBoardCopy()
//
//		// increment half move counter
//		ply++
//
//		// check that move is legal
//		if makeMove(moveList[count]) == 0 {
//			// decrement half move counter if move is illegal
//			ply--
//			continue
//		}
//
//		// increment legalMoves counter
//		legalMoves++
//
//		// score current move
//		score := -negaMax(-beta, -alpha, depth-1)
//
//		// decrement ply
//		ply--
//
//		// take move back
//		restoreBoardFromCopy(currentState)
//
//		// fail-hard beta cutoff
//		if score >= beta {
//			// node (move) fails high
//			return beta
//		}
//
//		// found a better move
//		if score > alpha {
//			// PV node (move)
//			alpha = score
//
//			// root node
//			if ply == 0 {
//				bestMoveSoFar = moveList[count]
//			}
//		}
//	}
//
//	// we don't have any legal moves to make in the current postion
//	if legalMoves == 0 {
//		// king is in check
//		if inCheck != 0 {
//			// debug
//			fmt.Println("checkmate found")
//			printBoard()
//
//			// return mating score (assuming closest distance to mating position)
//			return -49000 + ply
//		} else { // king not in check
//			// debug
//			printBoard()
//			fmt.Println("stalemate found")
//
//			// return stalemate score
//			return 0
//		}
//	}
//
//	// update alpha
//	if oldAlpha != alpha {
//		bestMove = bestMoveSoFar
//	}
//	// node (move) fails low
//	return alpha
//}

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

	// best move debug
	//fmt.Printf("best move: ")
	//if bestMove != 0 {
	//	printMove(bestMove)
	//}
	//fmt.Printf("\n")
	return bestMove
}
