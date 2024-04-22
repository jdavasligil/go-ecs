package bitset

type Bitset interface {
	GetBit(i int) bool
	SetBit(i int)
	ClearBit(i int)
	Len() int
	Reset()
}

type BitsetUint64 []uint64

func NewUint64(n int) BitsetUint64 {
	return make(BitsetUint64, (n+63)/64)
}

func (b BitsetUint64) GetBit(i int) bool {
	return (b[i/64] & (uint64(1) << (i % 64))) != 0
}

func (b BitsetUint64) SetBit(i int) {
	b[i/64] |= (uint64(1) << i % 64)
}

func (b BitsetUint64) ClearBit(i int) {
	b[i/64] &= ^(uint64(1) << i % 64)
}

func (b BitsetUint64) Len() int {
	return 64 * len(b)
}

func (b BitsetUint64) Reset() {
	for i := 0; i < len(b); i++ {
		b[i] = 0
	}
}
