package main

//*********************************
//           globals
//*********************************

// attacks
var bishopAttacks [64][512]uint64
var rookAttacks [64][4096]uint64
var pawnAttacks [2][64]uint64
var knightAttacks [64]uint64
var kingAttacks [64]uint64

// masks for generating sliding piece moves
var rookMasks [64]uint64
var bishopMasks [64]uint64

//*********************************
//            init
//*********************************

// loop through all the squares on the board generating attacks for pawns, knights, and kings on those squares
func initLeapersAttacks() {
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			square := uint64(rank*8 + file)
			pawnAttacks[white][square] = maskPawnAttacks(square, white)
			pawnAttacks[black][square] = maskPawnAttacks(square, black)
			knightAttacks[square] = maskKnightAttacks(square)
			kingAttacks[square] = maskKingAttacks(square)
		}
	}
}

// loop through all the squares on the board generating attacks for bishops and rooks
func initSlidersAttacks(piece int) {
	var attackMask uint64

	// loop over 64 squares and populate arrays with masks
	for square := uint64(0); square < 64; square++ {
		// init attack masks
		bishopMasks[square] = maskBishopAttacks(square)
		rookMasks[square] = maskRookAttacks(square)

		// init current mask
		if piece == bishop {
			// bishop
			attackMask = bishopMasks[square]
		} else if piece == rook {
			// rook
			attackMask = rookMasks[square]
		}

		// init relevant occupancy bit count
		relevantBitsCount := countBits(attackMask)

		// init occupancyIndices
		occupancyIndices := uint64(1) << relevantBitsCount

		// loop over occupancy indices
		for i := uint64(0); i < occupancyIndices; i++ {
			if piece == bishop {
				// bishop current occupancy
				occupancy := setOccupancy(i, relevantBitsCount, attackMask)

				// magic index
				magicIndex := (occupancy * bishopMagicNumbers[square]) >> (64 - bishopRelevantBits[square])

				// init current bishop attacks
				bishopAttacks[square][magicIndex] = bishopAttacksOnTheFly(square, occupancy)
			} else if piece == rook {
				// rook current occupancy
				occupancy := setOccupancy(i, relevantBitsCount, attackMask)
				//if square == h8 {
				//	printBitboard(setOccupancy(i, relevantBitsCount, attackMask))
				//	fmt.Println(i, relevantBitsCount, attackMask)
				//}

				// magic index
				magicIndex := (occupancy * rookMagicNumbers[square]) >> (64 - rookRelevantBits[square])

				// init current rook attacks
				rookAttacks[square][magicIndex] = rookAttacksOnTheFly(square, occupancy)
			}
		}
	}
}

//*********************************
//            leapers
//*********************************

// generate pawn attacks for a single square
func maskPawnAttacks(square uint64, side int) uint64 {
	var attacks uint64
	var leftAttack uint64
	var rightAttack uint64

	// set the pawn on an empty bitboard
	bitboard := setBit(0, square)

	// white pawns move up the board, black pawns move down
	// file masks prevent off board attacks
	if side == white {
		leftAttack = bitboard >> 9 & notHFile
		rightAttack = bitboard >> 7 & notAFile
	} else {
		leftAttack = bitboard << 7 & notHFile
		rightAttack = bitboard << 9 & notAFile
	}

	// merge bitboards
	attacks |= leftAttack | rightAttack
	return attacks
}

// generate knight attacks for a single square
func maskKnightAttacks(square uint64) uint64 {
	var attacks uint64

	// set the knight on an empty bitboard
	bitboard := setBit(0, square)

	// add offsets to attacks, leaving out off board attacks
	attacks |= bitboard >> 15 & notAFile
	attacks |= bitboard >> 6 & notABFile
	attacks |= bitboard << 10 & notABFile
	attacks |= bitboard << 17 & notAFile
	attacks |= bitboard << 15 & notHFile
	attacks |= bitboard << 6 & notGHFile
	attacks |= bitboard >> 10 & notGHFile
	attacks |= bitboard >> 17 & notHFile

	return attacks
}

// generate king attacks for a single square
func maskKingAttacks(square uint64) uint64 {
	var attacks uint64

	// set the king on an empty bitboard
	bitboard := setBit(0, square)

	// on the on board attacks, starting from 12 o clock and moving clockwise
	attacks |= bitboard >> 8
	attacks |= bitboard >> 7 & notAFile
	attacks |= bitboard << 1 & notAFile
	attacks |= bitboard << 9 & notAFile
	attacks |= bitboard << 8
	attacks |= bitboard << 7 & notHFile
	attacks |= bitboard >> 1 & notHFile
	attacks |= bitboard >> 9 & notHFile

	return attacks
}

//*********************************
//            bishop
//*********************************

// generate bishop inner attack mask for a single square for magic bitboard
func maskBishopAttacks(square uint64) uint64 {
	var attacks uint64

	// current rank and file
	currentRank := square / 8
	currentFile := square % 8

	// loop through squares in the four directions, stopping before hitting the edge of the board
	for targetRank, targetFile := currentRank+1, currentFile+1; targetRank < 7 && targetFile < 7; targetRank, targetFile = targetRank+1, targetFile+1 {
		targetSquare := targetRank*8 + targetFile
		attacks |= uint64(1) << targetSquare
	}

	for targetRank, targetFile := currentRank-1, currentFile+1; targetRank > 0 && targetFile < 7; targetRank, targetFile = targetRank-1, targetFile+1 {
		targetSquare := targetRank*8 + targetFile
		attacks |= uint64(1) << targetSquare
	}

	for targetRank, targetFile := currentRank-1, currentFile-1; targetRank > 0 && targetFile < 7 && targetFile > 0; targetRank, targetFile = targetRank-1, targetFile-1 {
		targetSquare := targetRank*8 + targetFile
		attacks |= uint64(1) << targetSquare
	}

	for targetRank, targetFile := currentRank+1, currentFile-1; targetRank < 7 && targetFile < 7 && targetFile > 0; targetRank, targetFile = targetRank+1, targetFile-1 {
		targetSquare := targetRank*8 + targetFile
		attacks |= uint64(1) << targetSquare
	}
	return attacks
}

// generate bishop attacks for a single square on the fly
func bishopAttacksOnTheFly(square uint64, blockers uint64) uint64 {
	var attacks uint64

	// current rank and file
	currentRank := square / 8
	currentFile := square % 8

	// loop through squares in the four directions, stopping after you hit a blocker
	for targetRank, targetFile := currentRank+1, currentFile+1; targetRank <= 7 && targetFile <= 7; targetRank, targetFile = targetRank+1, targetFile+1 {
		targetSquare := targetRank*8 + targetFile
		targetBitboard := uint64(1) << targetSquare
		attacks |= targetBitboard
		if (targetBitboard & blockers) != 0 {
			break
		}
	}

	for targetRank, targetFile := currentRank-1, currentFile+1; targetRank >= 0 && targetFile <= 7; targetRank, targetFile = targetRank-1, targetFile+1 {
		targetSquare := targetRank*8 + targetFile
		targetBitboard := uint64(1) << targetSquare
		attacks |= targetBitboard
		if (targetBitboard & blockers) != 0 {
			break
		}
	}

	for targetRank, targetFile := currentRank-1, currentFile-1; targetRank >= 0 && targetFile <= 7 && targetFile >= 0; targetRank, targetFile = targetRank-1, targetFile-1 {
		targetSquare := targetRank*8 + targetFile
		targetBitboard := uint64(1) << targetSquare
		attacks |= targetBitboard
		if (targetBitboard & blockers) != 0 {
			break
		}
	}

	for targetRank, targetFile := currentRank+1, currentFile-1; targetRank <= 7 && targetFile <= 7 && targetFile >= 0; targetRank, targetFile = targetRank+1, targetFile-1 {
		targetSquare := targetRank*8 + targetFile
		targetBitboard := uint64(1) << targetSquare
		attacks |= targetBitboard
		if (targetBitboard & blockers) != 0 {
			break
		}
	}
	return attacks
}

// get bishop attacks
func getBishopAttacks(square uint64, occupancy uint64) uint64 {
	// current board occupancy
	occupancy &= bishopMasks[square]
	occupancy *= bishopMagicNumbers[square]
	occupancy >>= 64 - bishopRelevantBits[square]

	// return attacks for relevant occupancy
	return bishopAttacks[square][occupancy]
}

//*********************************
//             rook
//*********************************

// generate rook inner attack mask for a single square for magic bitboard
func maskRookAttacks(square uint64) uint64 {
	var attacks uint64

	// current rank and file
	currentRank := square / 8
	currentFile := square % 8

	// loop through squares in the four directions, stopping before hitting the edge of the board
	for targetRank := currentRank + 1; targetRank < 7; targetRank++ {
		targetSquare := targetRank*8 + currentFile
		attacks |= uint64(1) << targetSquare

	}
	for targetRank := currentRank - 1; targetRank > 0 && targetRank < 7; targetRank-- {
		targetSquare := targetRank*8 + currentFile
		attacks |= uint64(1) << targetSquare

	}

	for targetFile := currentFile - 1; targetFile > 0 && targetFile < 7; targetFile-- {
		targetSquare := currentRank*8 + targetFile
		attacks |= uint64(1) << targetSquare

	}

	for targetFile := currentFile + 1; targetFile < 7; targetFile++ {
		targetSquare := currentRank*8 + targetFile
		attacks |= uint64(1) << targetSquare
	}
	return attacks
}

// generate rook attacks for a single square oin the fly
func rookAttacksOnTheFly(square uint64, blockers uint64) uint64 {
	var attacks uint64

	// current rank and file
	currentRank := square / 8
	currentFile := square % 8

	// loop through squares in the four directions, stopping after hitting a blocker
	for targetRank := currentRank + 1; targetRank <= 7; targetRank++ {
		targetSquare := targetRank*8 + currentFile
		targetBitboard := uint64(1) << targetSquare
		attacks |= targetBitboard
		if (targetBitboard & blockers) != 0 {
			break
		}
	}

	for targetRank := currentRank - 1; targetRank >= 0 && targetRank <= 7; targetRank-- {
		targetSquare := targetRank*8 + currentFile
		targetBitboard := uint64(1) << targetSquare
		attacks |= targetBitboard
		if (targetBitboard & blockers) != 0 {
			break
		}
	}

	for targetFile := currentFile - 1; targetFile >= 0 && targetFile <= 7; targetFile-- {
		targetSquare := currentRank*8 + targetFile
		targetBitboard := uint64(1) << targetSquare
		attacks |= targetBitboard
		if (targetBitboard & blockers) != 0 {
			break
		}
	}

	for targetFile := currentFile + 1; targetFile <= 7; targetFile++ {
		targetSquare := currentRank*8 + targetFile
		targetBitboard := uint64(1) << targetSquare
		attacks |= targetBitboard
		if (targetBitboard & blockers) != 0 {
			break
		}
	}
	return attacks
}

// get rook attacks
func getRookAttacks(square uint64, occupancy uint64) uint64 {
	// current board occupancy
	occupancy &= rookMasks[square]
	occupancy *= rookMagicNumbers[square]
	occupancy >>= 64 - rookRelevantBits[square]

	// return attacks for current occupancy
	return rookAttacks[square][occupancy]
}

//*********************************
//            queen
//*********************************

// queen moves are just rook + bishop moves
func getQueenAttacks(square uint64, occupancy uint64) uint64 {
	// current rook and bishop attacks
	queenAttacks := getRookAttacks(square, occupancy) | getBishopAttacks(square, occupancy)

	// return attacks for current occupancy
	return queenAttacks
}

//*********************************
//              misc
//*********************************

// set relevant occupancy bits
func setOccupancy(index uint64, bitsInMask uint64, attackMask uint64) uint64 {
	var occupancy uint64

	// loop over relevant bits
	for count := uint64(0); count < bitsInMask; count++ {
		// remove least significant bit from attack mask
		square := getLeastSignificantBitIndex(attackMask)
		attackMask = popBit(attackMask, square)
		// update occupancy if on the board
		if (index & (uint64(1) << count)) != 0 {
			occupancy |= uint64(1) << square
		}
	}
	return occupancy
}

// determine if a square is attacked
func isSquareAttacked(square uint64, side int) int {
	// attacked by black pawns
	if side == white { // white's turn
		if pawnAttacks[black][square]&bitboards[P] != 0 {
			return 1
		}
	} else if side == black { // attacked by white pawns
		if pawnAttacks[white][square]&bitboards[p] != 0 {
			return 1
		}
	}

	// attacked by knight
	if side == white {
		if knightAttacks[square]&bitboards[N] != 0 {
			return 1
		}
	} else if side == black {
		if knightAttacks[square]&bitboards[n] != 0 {
			return 1
		}
	}

	// attacked by king
	if side == white {
		if kingAttacks[square]&bitboards[K] != 0 {
			return 1
		}
	} else if side == black {
		if kingAttacks[square]&bitboards[k] != 0 {
			return 1
		}
	}

	// attacked by bishop
	if side == white {
		if getBishopAttacks(square, occupancies[both])&bitboards[B] != 0 {
			return 1
		}
	} else if side == black {
		if getBishopAttacks(square, occupancies[both])&bitboards[b] != 0 {
			return 1
		}
	}

	// attacked by rook
	if side == white {
		if getRookAttacks(square, occupancies[both])&bitboards[R] != 0 {
			return 1
		}
	} else if side == black {
		if getRookAttacks(square, occupancies[both])&bitboards[r] != 0 {
			return 1
		}
	}

	// attacked by queen
	if side == white {
		if getQueenAttacks(square, occupancies[both])&bitboards[Q] != 0 {
			return 1
		}
	} else if side == black {
		if getQueenAttacks(square, occupancies[both])&bitboards[q] != 0 {
			return 1
		}
	}

	return 0
}
