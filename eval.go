package main

func evaluate() int {
	score := 0
	for piece := P; piece <= k; piece++ {
		bitboard := bitboards[piece]
		for bitboard != 0 {
			// init square
			square := getLeastSignificantBitIndex(bitboard)

			// score material weights
			score += materialScore[piece]

			// score positional piece scores
			switch piece {
			// evaluate white piecesStr
			case P:
				score += pawnScore[square]
				break
			case N:
				score += knightScore[square]
				break
			case B:
				score += bishopScore[square]
				break
			case R:
				score += rookScore[square]
				break
			case K:
				score += kingScore[square]
				break

			// evaluate black piecesStr
			case p:
				score -= pawnScore[mirrorScore[square]]
				break
			case n:
				score -= knightScore[mirrorScore[square]]
				break
			case b:
				score -= bishopScore[mirrorScore[square]]
				break
			case r:
				score -= rookScore[mirrorScore[square]]
				break
			case k:
				score -= kingScore[mirrorScore[square]]
				break
			}

			// remove piece from bitboard
			bitboard = popBit(bitboard, square)
		}
	}
	if side == white {
		return score
	} else {
		return -score
	}
}
