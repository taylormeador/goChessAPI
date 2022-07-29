package main

//*********************************
//            bits
//*********************************

// check if a bit is on or off
func getBit(bitboard uint64, square uint64) uint64 {
	return bitboard & (uint64(1) << square)
}

// turn on a bit
func setBit(bitboard uint64, square uint64) uint64 {
	return bitboard | (uint64(1) << square)
}

// turn off a bit
func popBit(bitboard uint64, square uint64) uint64 {
	if bitboard == 1 {
		return bitboard & ^uint64(1)
	}
	return bitboard & ^(uint64(1) << square)
}

// count the number of bits on a bitboard
func countBits(bitboard uint64) uint64 {
	if bitboard == 0 {
		return 0
	}
	count := uint64(0)
	for {
		bitboard &= bitboard - 1
		count += 1
		if bitboard == 0 {
			break
		}
	}
	return count
}

// get the index of the least significant bit that is on
func getLeastSignificantBitIndex(bitboard uint64) uint64 {
	// check for a non empty bitboard
	if bitboard > 0 {
		leastSignificantBit := bitboard & -bitboard
		leadingOnes := leastSignificantBit - 1
		return countBits(leadingOnes)
	} else {
		// return out of range index if the board is empty
		return uint64(65)
	}

}

//*********************************
//            random
//*********************************

// generate pseudo random 32 bit number
var currentRandom = uint32(1804289383)

func getRandom32BitNumber() uint32 {
	// this number comes from Chess Programming's YouTube video
	number := currentRandom

	// XOR shift 32
	number ^= number << 13
	number ^= number >> 17
	number ^= number << 5

	// update current random
	currentRandom = number

	return number
}

// generate pseudo random 64 bit number
func getRandom64BitNumber() uint64 {
	var random64BitNumber uint64

	// get 4 random 32 bit numbers
	num1 := uint64(getRandom32BitNumber()) & 0xFFFF
	num2 := uint64(getRandom32BitNumber()) & 0xFFFF
	num3 := uint64(getRandom32BitNumber()) & 0xFFFF
	num4 := uint64(getRandom32BitNumber()) & 0xFFFF

	// shift bits
	random64BitNumber = num1 | (num2 << 16) | (num3 << 32) | (num4 << 48)

	return random64BitNumber
}
