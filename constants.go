package main

import "strings"

// enumerate colors
const (
	white = iota
	black
	both
)

func colorFromInt(side int) string {
	switch side {
	case 0:
		return "white"
	case 1:
		return "black"
	case 2:
		return "both"
	}
	return "error getting color"
}

// rook and bishop
const (
	rook   = 0
	bishop = 1
)

// enumerate board squares
const (
	a8 = uint64(iota)
	b8
	c8
	d8
	e8
	f8
	g8
	h8
	a7
	b7
	c7
	d7
	e7
	f7
	g7
	h7
	a6
	b6
	c6
	d6
	e6
	f6
	g6
	h6
	a5
	b5
	c5
	d5
	e5
	f5
	g5
	h5
	a4
	b4
	c4
	d4
	e4
	f4
	g4
	h4
	a3
	b3
	c3
	d3
	e3
	f3
	g3
	h3
	a2
	b2
	c2
	d2
	e2
	f2
	g2
	h2
	a1
	b1
	c1
	d1
	e1
	f1
	g1
	h1
	noSquare
)

// convert square name string to constant square uint64
func squareStringToUint64(s string) uint64 {
	switch s {
	case "a8":
		return a8
	case "b8":
		return b8
	case "c8":
		return c8
	case "d8":
		return d8
	case "e8":
		return e8
	case "f8":
		return f8
	case "g8":
		return g8
	case "h8":
		return h8
	case "a7":
		return a7
	case "b7":
		return b7
	case "c7":
		return c7
	case "d7":
		return d7
	case "e7":
		return e7
	case "f7":
		return f7
	case "g7":
		return g7
	case "h7":
		return h7
	case "a6":
		return a6
	case "b6":
		return b6
	case "c6":
		return c6
	case "d6":
		return d6
	case "e6":
		return e6
	case "f6":
		return f6
	case "g6":
		return g6
	case "h6":
		return h6
	case "a5":
		return a5
	case "b5":
		return b5
	case "c5":
		return c5
	case "d5":
		return d5
	case "e5":
		return e5
	case "f5":
		return f5
	case "g5":
		return g5
	case "h5":
		return h5
	case "a4":
		return a4
	case "b4":
		return b4
	case "c4":
		return c4
	case "d4":
		return d4
	case "e4":
		return e4
	case "f4":
		return f4
	case "g4":
		return g4
	case "h4":
		return h4
	case "a3":
		return a3
	case "b3":
		return b3
	case "c3":
		return c3
	case "d3":
		return d3
	case "e3":
		return e3
	case "f3":
		return f3
	case "g3":
		return g3
	case "h3":
		return h3
	case "a2":
		return a2
	case "b2":
		return b2
	case "c2":
		return c2
	case "d2":
		return d2
	case "e2":
		return e2
	case "f2":
		return f2
	case "g2":
		return g2
	case "h2":
		return h2
	case "a1":
		return a1
	case "b1":
		return b1
	case "c1":
		return c1
	case "d1":
		return d1
	case "e1":
		return e1
	case "f1":
		return f1
	case "g1":
		return g1
	case "h1":
		return h1
	}
	return noSquare
}

// enumerate pieces
// uppercase represents white pieces, lowercase represents black
const (
	P = iota
	N
	B
	R
	Q
	K
	p
	n
	b
	r
	q
	k
)

// bitboard masks
const (
	//  e.g. all 0's in the "a" file, all 1's elsewhere
	//
	//  8    0  1  1  1  1  1  1  1
	//  7    0  1  1  1  1  1  1  1
	//  6    0  1  1  1  1  1  1  1
	//  5    0  1  1  1  1  1  1  1
	//  4    0  1  1  1  1  1  1  1
	//  3    0  1  1  1  1  1  1  1
	//  2    0  1  1  1  1  1  1  1
	//  1    0  1  1  1  1  1  1  1
	//
	//       a  b  c  d  e  f  g  h
	//       Bitboard: 18374403900871474942

	notAFile    = uint64(18374403900871474942)
	notHFile    = uint64(9187201950435737471)
	notABFile   = uint64(18229723555195321596)
	notGHFile   = uint64(4557430888798830399)
	firstRank   = uint64(18302628885633695745)
	secondRank  = uint64(71776119061217280)
	seventhRank = uint64(65280)
	eighthRank  = uint64(255)
)

// castling rights
const (
	// bq bk wq wk
	// 1  1  1  1   => all castling rights
	// 1000 => black can castle queen side
	// 1101 => white cannot castle queen side
	// 0000 => no one can castle
	wk = 1
	wq = 2
	bk = 4
	bq = 8
)

var pieceToCastle = map[byte]int{
	'K': wk, 'Q': wq, 'k': bk, 'q': bq,
}

func castleToString(castle int) string {
	var sb strings.Builder

	if castle&wk != 0 {
		sb.WriteString("K")
	}
	if castle&wq != 0 {
		sb.WriteString("Q")
	}
	if castle&bk != 0 {
		sb.WriteString("k")
	}
	if castle&bq != 0 {
		sb.WriteString("q")
	}
	return sb.String()
}

// string pieces
var stringPieces = [12]string{
	"P", "N", "B", "R", "Q", "K", "p", "n", "b", "r", "q", "k",
}

// ASCII pieces
var asciiPieces = [12]byte{
	'P', 'N', 'B', 'R', 'Q', 'K', 'p', 'n', 'b', 'r', 'q', 'k',
}

// unicode pieces
var unicodePieces = [12]string{
	"♟", "♞", "♝", "♜", "♛", "♚", "♙", "♘", "♗", "♖", "♕", "♔",
}

// convert ASCII character to encoded constants
var charPieces = map[byte]int{
	'P': 0, 'N': 1, 'B': 2, 'R': 3, 'Q': 4, 'K': 5, 'p': 6, 'n': 7, 'b': 8, 'r': 9, 'q': 10, 'k': 11,
}

// convert promoted piece move encoding to string piece
var promotedPieces = map[uint64]string{
	0:  " ",
	1:  "N",
	2:  "B",
	3:  "R",
	4:  "Q",
	5:  " ",
	6:  " ",
	7:  "n",
	8:  "b",
	9:  "r",
	10: "q",
	11: " ",
}

// lookup name of square from index in bitboard
var algebraic = [65]string{
	"a8", "b8", "c8", "d8", "e8", "f8", "g8", "h8",
	"a7", "b7", "c7", "d7", "e7", "f7", "g7", "h7",
	"a6", "b6", "c6", "d6", "e6", "f6", "g6", "h6",
	"a5", "b5", "c5", "d5", "e5", "f5", "g5", "h5",
	"a4", "b4", "c4", "d4", "e4", "f4", "g4", "h4",
	"a3", "b3", "c3", "d3", "e3", "f3", "g3", "h3",
	"a2", "b2", "c2", "d2", "e2", "f2", "g2", "h2",
	"a1", "b1", "c1", "d1", "e1", "f1", "g1", "h1", "-",
}

// relevant occupancy bit counts for a rook at every square on board
var rookRelevantBits = [64]uint64{
	12, 11, 11, 11, 11, 11, 11, 12,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	12, 11, 11, 11, 11, 11, 11, 12,
}

// relevant occupancy bit counts for a bishop at every square on board
var bishopRelevantBits = [64]uint64{
	6, 5, 5, 5, 5, 5, 5, 6,
	5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5,
	6, 5, 5, 5, 5, 5, 5, 6,
}

// rook magic numbers from initMagicNumbers()
var rookMagicNumbers = [64]uint64{
	0x8a80104000800020,
	0x140002000100040,
	0x2801880a0017001,
	0x100081001000420,
	0x200020010080420,
	0x3001c0002010008,
	0x8480008002000100,
	0x2080088004402900,
	0x800098204000,
	0x2024401000200040,
	0x100802000801000,
	0x120800800801000,
	0x208808088000400,
	0x2802200800400,
	0x2200800100020080,
	0x801000060821100,
	0x80044006422000,
	0x100808020004000,
	0x12108a0010204200,
	0x140848010000802,
	0x481828014002800,
	0x8094004002004100,
	0x4010040010010802,
	0x20008806104,
	0x100400080208000,
	0x2040002120081000,
	0x21200680100081,
	0x20100080080080,
	0x2000a00200410,
	0x20080800400,
	0x80088400100102,
	0x80004600042881,
	0x4040008040800020,
	0x440003000200801,
	0x4200011004500,
	0x188020010100100,
	0x14800401802800,
	0x2080040080800200,
	0x124080204001001,
	0x200046502000484,
	0x480400080088020,
	0x1000422010034000,
	0x30200100110040,
	0x100021010009,
	0x2002080100110004,
	0x202008004008002,
	0x20020004010100,
	0x2048440040820001,
	0x101002200408200,
	0x40802000401080,
	0x4008142004410100,
	0x2060820c0120200,
	0x1001004080100,
	0x20c020080040080,
	0x2935610830022400,
	0x44440041009200,
	0x280001040802101,
	0x2100190040002085,
	0x80c0084100102001,
	0x4024081001000421,
	0x20030a0244872,
	0x12001008414402,
	0x2006104900a0804,
	0x1004081002402,
}

// bishop magic numbers from initMagicNumbers()
var bishopMagicNumbers = [64]uint64{
	0x40040844404084,
	0x2004208a004208,
	0x10190041080202,
	0x108060845042010,
	0x581104180800210,
	0x2112080446200010,
	0x1080820820060210,
	0x3c0808410220200,
	0x4050404440404,
	0x21001420088,
	0x24d0080801082102,
	0x1020a0a020400,
	0x40308200402,
	0x4011002100800,
	0x401484104104005,
	0x801010402020200,
	0x400210c3880100,
	0x404022024108200,
	0x810018200204102,
	0x4002801a02003,
	0x85040820080400,
	0x810102c808880400,
	0xe900410884800,
	0x8002020480840102,
	0x220200865090201,
	0x2010100a02021202,
	0x152048408022401,
	0x20080002081110,
	0x4001001021004000,
	0x800040400a011002,
	0xe4004081011002,
	0x1c004001012080,
	0x8004200962a00220,
	0x8422100208500202,
	0x2000402200300c08,
	0x8646020080080080,
	0x80020a0200100808,
	0x2010004880111000,
	0x623000a080011400,
	0x42008c0340209202,
	0x209188240001000,
	0x400408a884001800,
	0x110400a6080400,
	0x1840060a44020800,
	0x90080104000041,
	0x201011000808101,
	0x1a2208080504f080,
	0x8012020600211212,
	0x500861011240000,
	0x180806108200800,
	0x4000020e01040044,
	0x300000261044000a,
	0x802241102020002,
	0x20906061210001,
	0x5a84841004010310,
	0x4010801011c04,
	0xa010109502200,
	0x4a02012000,
	0x500201010098b028,
	0x8040002811040900,
	0x28000010020204,
	0x6000020202d0240,
	0x8918844842082200,
	0x4010011029020020,
}

// FEN debug strings
const (
	emptyBoard        = "8/8/8/8/8/8/8/8 w - - 0 1"
	startPosition     = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	trickyPosition    = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"
	killerPosition    = "rnbqkb1r/pp1p1pPp/8/2p1pP2/1P1P4/3P3P/P1P1P3/RNBQKBNR w KQkq e6 0 1"
	cmkPosition       = "r2q1rk1/ppp2ppp/2n1bn2/2b1p3/3pP3/3P1NPP/PPP1NPB1/R1BQ1RK1 b - - 0 9"
	checkmatePosition = "rnbqkbnr/ppppp2p/5p2/8/8/4P3/PPPP1PPP/RNBQKBNR w KQkq - 0 1"
	stalematePosition = "4k3/R7/8/8/8/8/8/3R1R1K b - - 0 1"
)

type FENjson struct {
	FEN   string `json:"FEN"`
	legal bool   `json:"legal"`
}

type bestMovejson struct {
	FEN   string `json:"FEN"`
	legal bool   `json:"legal"`
	best  string `json:"best"`
}

/*
                           castling   move     in      in
                              right update     binary  decimal
 king & rooks didn't move:     1111 & 1111  =  1111    15
        white king  moved:     1111 & 1100  =  1100    12
  white king's rook moved:     1111 & 1110  =  1110    14
 white queen's rook moved:     1111 & 1101  =  1101    13

         black king moved:     1111 & 0011  =  1011    3
  black king's rook moved:     1111 & 1011  =  1011    11
 black queen's rook moved:     1111 & 0111  =  0111    7
*/
var castlingRights = [64]int{
	7, 15, 15, 15, 3, 15, 15, 11,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	13, 15, 15, 15, 12, 15, 15, 14,
}

// material scrore
/*
   ♙ =   100   = ♙
   ♘ =   300   = ♙ * 3
   ♗ =   350   = ♙ * 3 + ♙ * 0.5
   ♖ =   500   = ♙ * 5
   ♕ =   1000  = ♙ * 10
   ♔ =   10000 = ♙ * 100

*/
var materialScore = [12]int{
	100,    // white pawn score
	300,    // white knight scrore
	350,    // white bishop score
	500,    // white rook score
	1000,   // white queen score
	10000,  // white king score
	-100,   // black pawn score
	-300,   // black knight scrore
	-350,   // black bishop score
	-500,   // black rook score
	-1000,  // black queen score
	-10000, // black king score
}

// pawn positional score
var pawnScore = [64]int{
	90, 90, 90, 90, 90, 90, 90, 90,
	30, 30, 30, 40, 40, 30, 30, 30,
	20, 20, 20, 30, 30, 30, 20, 20,
	10, 10, 10, 20, 20, 10, 10, 10,
	5, 5, 10, 20, 20, 5, 5, 5,
	0, 0, 0, 5, 5, 0, 0, 0,
	0, 0, 0, -10, -10, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

// knight positional score
var knightScore = [64]int{
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 10, 10, 0, 0, -5,
	-5, 5, 20, 20, 20, 20, 5, -5,
	-5, 10, 20, 30, 30, 20, 10, -5,
	-5, 10, 20, 30, 30, 20, 10, -5,
	-5, 5, 20, 10, 10, 20, 5, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, -10, 0, 0, 0, 0, -10, -5,
}

// bishop positional score
var bishopScore = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 10, 10, 0, 0, 0,
	0, 0, 10, 20, 20, 10, 0, 0,
	0, 0, 10, 20, 20, 10, 0, 0,
	0, 10, 0, 0, 0, 0, 10, 0,
	0, 30, 0, 0, 0, 0, 30, 0,
	0, 0, -10, 0, 0, -10, 0, 0,
}

// rook positional score
var rookScore = [64]int{
	50, 50, 50, 50, 50, 50, 50, 50,
	50, 50, 50, 50, 50, 50, 50, 50,
	0, 0, 10, 20, 20, 10, 0, 0,
	0, 0, 10, 20, 20, 10, 0, 0,
	0, 0, 10, 20, 20, 10, 0, 0,
	0, 0, 10, 20, 20, 10, 0, 0,
	0, 0, 10, 20, 20, 10, 0, 0,
	0, 0, 0, 20, 20, 0, 0, 0,
}

// king positional score
var kingScore = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 5, 5, 5, 5, 0, 0,
	0, 5, 5, 10, 10, 5, 5, 0,
	0, 5, 10, 20, 20, 10, 5, 0,
	0, 5, 10, 20, 20, 10, 5, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 5, 5, -5, -5, 0, 5, 0,
	0, 0, 5, 0, -15, 0, 10, 0,
}

// mirror positional score tables for opposite side
var mirrorScore = [128]uint64{
	a1, b1, c1, d1, e1, f1, g1, h1,
	a2, b2, c2, d2, e2, f2, g2, h2,
	a3, b3, c3, d3, e3, f3, g3, h3,
	a4, b4, c4, d4, e4, f4, g4, h4,
	a5, b5, c5, d5, e5, f5, g5, h5,
	a6, b6, c6, d6, e6, f6, g6, h6,
	a7, b7, c7, d7, e7, f7, g7, h7,
	a8, b8, c8, d8, e8, f8, g8, h8,
}
