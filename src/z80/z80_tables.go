package z80

/* Whether a half carry occurred or not can be determined by looking at
   the 3rd bit of the two arguments and the result; these are hashed
   into this table in the form r12, where r is the 3rd bit of the
   result, 1 is the 3rd bit of the 1st argument and 2 is the
   third bit of the 2nd argument; the tables differ for add and subtract
   operations */
var halfcarryAddTable = []byte{0, FLAG_H, FLAG_H, FLAG_H, 0, 0, 0, FLAG_H}
var halfcarrySubTable = []byte{0, 0, FLAG_H, 0, FLAG_H, 0, FLAG_H, FLAG_H}

/* Similarly, overflow can be determined by looking at the 7th bits; again
   the hash into this table is r12 */
var overflowAddTable = []byte{0, 0, 0, FLAG_V, FLAG_V, 0, 0, 0}
var overflowSubTable = []byte{0, FLAG_V, 0, 0, 0, 0, FLAG_V, 0}

var sz53Table, sz53pTable, parityTable [0x100]byte

func init() {
	var i int16
	var j, k byte
	var parity byte

	for i = 0; i < 0x100; i++ {
		sz53Table[i] = byte(i) & (0x08 | 0x20 | 0x80)
		j = byte(i)
		parity = 0
		for k = 0; k < 8; k++ {
			parity ^= j & 1
			j >>= 1
		}
		parityTable[i] = ternOpB(parity != 0, 0, 0x04)
		sz53pTable[i] = sz53Table[i] | parityTable[i]
	}

	sz53Table[0] |= 0x40
	sz53pTable[0] |= 0x40
}
