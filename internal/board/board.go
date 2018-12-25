package board

import "sync"

type BoardReader interface {
	Get(i, j int) FillState
}

type Board interface {
	BoardReader
	Set(i, j int, s FillState) bool
}

type board struct {
	mu sync.RWMutex
	s  []FillState
}

func New() Board {
	return &board{
		s: make([]FillState, 9),
	}
}

type FillState int

func (s FillState) String() string {
	switch s {
	case FillStateX:
		return "X"
	case FillStateO:
		return "O"
	default:
		return " "
	}
}

const (
	FillStateEmpty FillState = iota
	FillStateX
	FillStateO
)

// Get returns the current state for the given coordinates. If anything beyond
// the limits is requested, the method will panic.
//
// +---------------------------+
// |i=0,j=0 | i=1,j=0 | i=2,j=0|
// +---------------------------+
// |i=0,j=1 | i=1,j=1 | i=2,j=1|
// +---------------------------+
// |i=0,j=2 | i=1,j=2 | i=2,j=2|
// +---------------------------+
func (b *board) Get(i, j int) FillState {
	b.checkBounds(i, j)
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.s[j*3+i]
}

// Set will set the state of the cooridate to the given value. If the value
// has already been filled, then the method will return false. Otherwise it
// will return true. If anything beyond the limits is tried, the method will
// panic.
func (b *board) Set(i, j int, s FillState) bool {
	b.checkBounds(i, j)
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.s[j*3+i] != FillStateEmpty {
		return false
	}
	b.s[j*3+i] = s
	return true
}

func (b *board) checkBounds(i, j int) {
	if i < 0 || j < 0 || i > 2 || j > 2 {
		panic("out of bounds")
	}
}
