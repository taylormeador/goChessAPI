package main

import "fmt"

// generate moves for a given square
func generateMoves(sourceSquare uint64, moveList *[]uint64) {
	startSquare := uint64(1) << sourceSquare

	var targetSquare, move uint64
	var promotionRank, promotionRankMinusOne, pawnStartRank, pawnPush, doublePawnPush, pawnPromotion uint64
	var pawnPushOffset, doublePawnPushOffset, kingStartSquare uint64
	var pawn, knight, bishop, rook, queen, king uint64
	var bSquare, cSquare, dSquare, fSquare, gSquare uint64
	var castleKingside, castleQueenside, enemyColor int

	// assign side related variables
	if side == white {
		pawn = P
		knight = N
		bishop = B
		rook = R
		queen = Q
		king = K
		enemyColor = black
		promotionRank = eighthRank
		promotionRankMinusOne = seventhRank
		pawnStartRank = secondRank
		pawnPush = startSquare >> 8
		pawnPushOffset = sourceSquare - 8
		doublePawnPush = (startSquare & pawnStartRank) >> 16
		doublePawnPushOffset = sourceSquare - 16
		pawnPromotion = (startSquare & promotionRankMinusOne) >> 8
		castleKingside = castle & wk
		castleQueenside = castle & wq
		kingStartSquare = e1
		bSquare = b1
		cSquare = c1
		dSquare = d1
		fSquare = f1
		gSquare = g1

	} else if side == black {
		pawn = p
		knight = n
		bishop = b
		rook = r
		queen = q
		king = k
		enemyColor = white
		promotionRank = firstRank
		promotionRankMinusOne = secondRank
		pawnStartRank = seventhRank
		pawnPush = startSquare << 8
		pawnPushOffset = sourceSquare + 8
		doublePawnPush = (startSquare & pawnStartRank) << 16
		doublePawnPushOffset = sourceSquare + 16
		pawnPromotion = (startSquare & promotionRankMinusOne) << 8
		castleKingside = castle & bk
		castleQueenside = castle & bq
		kingStartSquare = e8
		bSquare = b8
		cSquare = c8
		dSquare = d8
		fSquare = f8
		gSquare = g8
	}

	// quiet pawn moves
	if bitboards[pawn]&startSquare != 0 { // if there is a white pawn on the start square
		// single pawn push
		targetSquare = (pawnPush & ^occupancies[both]) & ^promotionRank
		if targetSquare != 0 {
			move = encodeMove(sourceSquare, pawnPushOffset, pawn, 0, 0, 0, 0, 0)
			addMove(move, moveList)
		}
		// double pawn push
		targetSquare = doublePawnPush & ^occupancies[both]
		squareInFront := pawnPush & ^occupancies[both]
		if targetSquare != 0 && squareInFront != 0 { // check that the squares in front and two squares in front are empty
			addMove(encodeMove(sourceSquare, doublePawnPushOffset, pawn, 0, 0, 1, 0, 0), moveList)
		}
		// pawn promotion
		targetSquare = pawnPromotion & ^occupancies[both]
		if targetSquare != 0 {
			addMove(encodeMove(sourceSquare, pawnPushOffset, pawn, queen, 0, 0, 0, 0), moveList)
			addMove(encodeMove(sourceSquare, pawnPushOffset, pawn, rook, 0, 0, 0, 0), moveList)
			addMove(encodeMove(sourceSquare, pawnPushOffset, pawn, bishop, 0, 0, 0, 0), moveList)
			addMove(encodeMove(sourceSquare, pawnPushOffset, pawn, knight, 0, 0, 0, 0), moveList)
		}
	}

	// pawn captures
	if bitboards[pawn]&startSquare != 0 {
		pawnCaptures := pawnAttacks[side][sourceSquare] & occupancies[enemyColor]
		for {
			if pawnCaptures != 0 {
				targetSquare = getLeastSignificantBitIndex(pawnCaptures)
				if startSquare&promotionRankMinusOne != 0 { // promotion capture
					addMove(encodeMove(sourceSquare, targetSquare, pawn, queen, 1, 0, 0, 0), moveList)
					addMove(encodeMove(sourceSquare, targetSquare, pawn, rook, 1, 0, 0, 0), moveList)
					addMove(encodeMove(sourceSquare, targetSquare, pawn, knight, 1, 0, 0, 0), moveList)
					addMove(encodeMove(sourceSquare, targetSquare, pawn, bishop, 1, 0, 0, 0), moveList)
				} else { // regular capture
					addMove(encodeMove(sourceSquare, targetSquare, pawn, 0, 1, 0, 0, 0), moveList)
				}
				pawnCaptures = popBit(pawnCaptures, targetSquare)
			} else {
				break
			}
		}
		// en passant
		if enPassantSquare != noSquare {
			enPassantCapture := pawnAttacks[side][sourceSquare] & (1 << enPassantSquare)
			if enPassantCapture != 0 {
				targetSquare = getLeastSignificantBitIndex(enPassantCapture)
				addMove(encodeMove(sourceSquare, targetSquare, pawn, 0, 1, 0, 1, 0), moveList)
			}
		}
	}

	// castling moves
	if sourceSquare == kingStartSquare { // white
		if castleKingside != 0 { // kingside
			if getBit(occupancies[both], fSquare) == 0 && getBit(occupancies[both], gSquare) == 0 { // check that the kingside squares are empty
				if isSquareAttacked(fSquare, enemyColor) == 0 && isSquareAttacked(kingStartSquare, enemyColor) == 0 { // check that travel square and the king are not attacked
					move = encodeMove(sourceSquare, gSquare, king, 0, 0, 0, 0, 1)
					addMove(move, moveList)
				}
			}
		}
		if castleQueenside != 0 { // queenside
			if getBit(occupancies[both], bSquare) == 0 && getBit(occupancies[both], cSquare) == 0 && getBit(occupancies[both], dSquare) == 0 { // check that the queenside squares are empty
				if isSquareAttacked(dSquare, enemyColor) == 0 && isSquareAttacked(kingStartSquare, enemyColor) == 0 { // check that travel squares and the king are not attacked
					move = encodeMove(sourceSquare, cSquare, king, 0, 0, 0, 0, 1)
					addMove(move, moveList)
				}
			}
		}
	}

	// workin on the Night Moves
	if bitboards[knight]&startSquare != 0 { // there is a knight on the sourceSquare
		// quiet knight moves
		quietMoves := knightAttacks[sourceSquare] &^ occupancies[both]
		for {
			if quietMoves != 0 { // knight moves to empty squares
				targetSquare = getLeastSignificantBitIndex(quietMoves)
				addMove(encodeMove(sourceSquare, targetSquare, knight, 0, 0, 0, 0, 0), moveList)
				quietMoves = popBit(quietMoves, targetSquare)
			} else {
				break
			}
		}
		// knight attacks
		attackMoves := knightAttacks[sourceSquare] & occupancies[enemyColor]
		for {
			if attackMoves != 0 { // the knight is attacking an enemy piece
				targetSquare = getLeastSignificantBitIndex(attackMoves)
				addMove(encodeMove(sourceSquare, targetSquare, knight, 0, 1, 0, 0, 0), moveList)
				attackMoves = popBit(attackMoves, targetSquare)
			} else {
				break
			}
		}
	}

	// bishop moves
	if bitboards[bishop]&startSquare != 0 { // if there's a bishop on the square
		bishopMoves := getBishopAttacks(sourceSquare, occupancies[both]) & ^occupancies[side]
		for {
			if bishopMoves != 0 {
				targetSquare = getLeastSignificantBitIndex(bishopMoves)
				if occupancies[enemyColor]&(1<<targetSquare) != 0 {
					addMove(encodeMove(sourceSquare, targetSquare, bishop, 0, 1, 0, 0, 0), moveList)
				} else {
					addMove(encodeMove(sourceSquare, targetSquare, bishop, 0, 0, 0, 0, 0), moveList)
				}
				bishopMoves = popBit(bishopMoves, targetSquare)
			} else {
				break
			}
		}
	}

	// rook moves
	if bitboards[rook]&startSquare != 0 { // there is a rook on the square
		rookMoves := getRookAttacks(sourceSquare, occupancies[both]) & ^occupancies[side]
		for {
			if rookMoves != 0 {
				targetSquare = getLeastSignificantBitIndex(rookMoves)
				if occupancies[enemyColor]&(1<<targetSquare) != 0 {
					addMove(encodeMove(sourceSquare, targetSquare, rook, 0, 1, 0, 0, 0), moveList)
				} else {
					addMove(encodeMove(sourceSquare, targetSquare, rook, 0, 0, 0, 0, 0), moveList)
				}
				rookMoves = popBit(rookMoves, targetSquare)
			} else {
				break
			}
		}
	}

	// queen moves
	if bitboards[queen]&startSquare != 0 { // it is white's turn and there is a white queen on the square
		queenMoves := getQueenAttacks(sourceSquare, occupancies[both]) & ^occupancies[side]
		for {
			if queenMoves != 0 {
				targetSquare = getLeastSignificantBitIndex(queenMoves)
				if occupancies[enemyColor]&(1<<targetSquare) != 0 { // queen attacking enemy piece
					addMove(encodeMove(sourceSquare, targetSquare, queen, 0, 1, 0, 0, 0), moveList)
				} else {
					addMove(encodeMove(sourceSquare, targetSquare, queen, 0, 0, 0, 0, 0), moveList)
				}
				queenMoves = popBit(queenMoves, targetSquare)
			} else {
				break
			}
		}
	}

	// king moves
	if bitboards[king]&startSquare != 0 {
		kingMoves := kingAttacks[sourceSquare] & ^occupancies[side]
		for {
			if kingMoves != 0 {
				targetSquare = getLeastSignificantBitIndex(kingMoves)
				if occupancies[enemyColor]&(1<<targetSquare) != 0 { // king attacking enemy piece
					addMove(encodeMove(sourceSquare, targetSquare, king, 0, 0, 0, 0, 0), moveList)
				} else {
					addMove(encodeMove(sourceSquare, targetSquare, king, 0, 0, 0, 0, 0), moveList)
				}
				kingMoves = popBit(kingMoves, targetSquare)
			} else {
				break
			}
		}
	}
}

// generate moves for all squares
func generateAllMoves(moveList []uint64) []uint64 {
	for rank := uint64(0); rank < 8; rank++ {
		for file := uint64(0); file < 8; file++ {
			square := rank*8 + file
			generateMoves(square, &moveList)
		}
	}
	return moveList
}

/*  encode moves in binary
      binary move bits                               hexadecimal constants

0000 0000 0000 0000 0011 1111    source square       0x3f
0000 0000 0000 1111 1100 0000    target square       0xfc0
0000 0000 1111 0000 0000 0000    piece               0xf000
0000 1111 0000 0000 0000 0000    promoted piece      0xf0000
0001 0000 0000 0000 0000 0000    capture flag        0x100000
0010 0000 0000 0000 0000 0000    double push flag    0x200000
0100 0000 0000 0000 0000 0000    enpassant flag      0x400000
1000 0000 0000 0000 0000 0000    castling flag       0x800000
*/
func encodeMove(source uint64, target uint64, piece uint64, promoted uint64, capture uint64,
	double uint64, enPassant uint64, castling uint64) uint64 {

	move := source | (target << 6) | (piece << 12) | (promoted << 16) | (capture << 20) |
		(double << 21) | (enPassant << 22) | (castling << 23)

	return move
}

// return the specified attribute of an encoded move
func getMoveAttr(move uint64, attr string) uint64 {
	switch attr {
	case "source":
		return move & 0x3f
	case "target": // target
		return (move & 0xfc0) >> 6
	case "piece": // piece
		return (move & 0xf000) >> 12
	case "promoted": // promoted
		return (move & 0xf0000) >> 16
	case "capture": // capture
		if move&0x100000 != 0 {
			return 1
		}
		return 0

	case "double": // double
		if move&0x200000 != 0 {
			return 1
		}
		return 0
	case "enPassant": // enPassant
		if move&0x400000 != 0 {
			return 1
		}
		return 0
	case "castling": // castling
		if move&0x800000 != 0 {
			return 1
		}
		return 0
	}
	return 33554431
}

// append move to moveList
func addMove(move uint64, moveList *[]uint64) {
	*moveList = append(*moveList, move)
}

// print move source, target, and promoted piece
func printMove(move uint64) {
	fmt.Printf("%s%s%s", algebraic[getMoveAttr(move, "source")],
		algebraic[getMoveAttr(move, "target")], promotedPieces[getMoveAttr(move, "promoted")],
	)
}

// loop through all moves in move list and print
func printMoveList(moveList []uint64) {
	// formatting
	fmt.Printf("\n    move    piece   capture   double    enpass    castling\n\n")

	// loop through movesList
	for _, move := range moveList {
		fmt.Printf("    %s%s%s   %s       %d         %d         %d         %d\n",
			algebraic[getMoveAttr(move, "source")],
			algebraic[getMoveAttr(move, "target")],
			promotedPieces[getMoveAttr(move, "promoted")],
			stringPieces[getMoveAttr(move, "piece")],
			getMoveAttr(move, "capture"),
			getMoveAttr(move, "double"),
			getMoveAttr(move, "enPassant"),
			getMoveAttr(move, "castling"),
		)
	}
	// print total number of moves
	fmt.Println()
	fmt.Printf("    Total number of moves: %d\n\n", len(moveList))
}

func makeMove(move uint64) int {
	// preserve state
	gameState := returnBoardCopy()

	// decode move
	source := getMoveAttr(move, "source")
	target := getMoveAttr(move, "target")
	promoted := getMoveAttr(move, "promoted")
	piece := getMoveAttr(move, "piece")
	capture := getMoveAttr(move, "capture")
	double := getMoveAttr(move, "double")
	enPassant := getMoveAttr(move, "enPassant")
	castling := getMoveAttr(move, "castling")

	// remove piece from source square
	bitboards[piece] = popBit(bitboards[piece], source)

	// place piece on target square
	bitboards[piece] = setBit(bitboards[piece], target)

	// handle promotion
	if promoted != 0 {
		// remove piece from target square
		bitboards[piece] = popBit(bitboards[piece], target)

		// place promoted piece on target square
		bitboards[promoted] = setBit(bitboards[promoted], target)
	}

	// handle capture
	if capture != 0 {
		// set piece bounds for loop
		var startPiece, endPiece uint64
		if side == white {
			startPiece = p
			endPiece = k
		} else if side == black {
			startPiece = P
			endPiece = K
		}

		// loop through enemy piece bitboards and remove the captured piece, if it's there
		for targetPiece := startPiece; targetPiece <= endPiece; targetPiece++ {
			if getBit(bitboards[targetPiece], target) != 0 {
				bitboards[targetPiece] = popBit(bitboards[targetPiece], target)
				break
			}
		}
	}

	// handle double pawn moves
	if double != 0 {
		// set enPassantSquare
		if side == white {
			enPassantSquare = target + 8
		} else if side == black {
			enPassantSquare = target - 8
		}
	} else {
		enPassantSquare = noSquare
	}

	// handle en passant
	if enPassant != 0 {
		// need to remove the piece one rank ahead or behind the target, depending on side
		if side == white {
			bitboards[p] = popBit(bitboards[p], target+8)
		} else if side == black {
			bitboards[P] = popBit(bitboards[P], target-8)
		}
	}

	// handle castling moves
	if castling != 0 {
		switch target {
		// white castles king side
		case g1:
			// move the h rook
			bitboards[R] = popBit(bitboards[R], h1)
			bitboards[R] = setBit(bitboards[R], f1)

		// white castles queen side
		case c1:
			// move the a rook
			bitboards[R] = popBit(bitboards[R], a1)
			bitboards[R] = setBit(bitboards[R], d1)

		// black castles king side
		case g8:
			// move the h rook
			bitboards[r] = popBit(bitboards[r], h8)
			bitboards[r] = setBit(bitboards[r], f8)

		// black castles queen side
		case c8:
			// move the a rook
			bitboards[r] = popBit(bitboards[r], a8)
			bitboards[r] = setBit(bitboards[r], d8)
		}
	}

	// update castling rights
	castle &= castlingRights[source]

	// reset then update occupancies
	occupancies[white], occupancies[black], occupancies[both] = 0, 0, 0

	// loop through white pieces
	for whitePiece := P; whitePiece <= K; whitePiece++ {
		occupancies[white] |= bitboards[whitePiece]
	}

	// loop through black pieces
	for blackPiece := p; blackPiece <= k; blackPiece++ {
		occupancies[black] |= bitboards[blackPiece]
	}

	// both occupancies
	occupancies[both] = occupancies[white] | occupancies[black]

	// ensure that the king has not moved into check
	kingLocation := a1
	enemySide := white
	if side == white {
		kingLocation = getLeastSignificantBitIndex(bitboards[K])
		enemySide = black
	} else {
		kingLocation = getLeastSignificantBitIndex(bitboards[k])
	}

	if isSquareAttacked(kingLocation, enemySide) != 0 {
		// revert the game state
		restoreBoardFromCopy(gameState)
		return 0
	} else {
		side = enemySide
		return 1
	}

}
