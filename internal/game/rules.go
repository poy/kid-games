package game

import (
	"github.com/poy/kids-games/internal/board"
)

func Winner(b board.Board) board.FillState {
	for _, s := range []struct {
		I []int
		J []int
	}{
		{
			I: []int{0, 0, 0},
			J: []int{0, 1, 2},
		},
		{
			I: []int{1, 1, 1},
			J: []int{0, 1, 2},
		},
		{
			I: []int{2, 2, 2},
			J: []int{0, 1, 2},
		},
		{
			I: []int{0, 1, 2},
			J: []int{0, 0, 0},
		},
		{
			I: []int{0, 1, 2},
			J: []int{1, 1, 1},
		},
		{
			I: []int{0, 1, 2},
			J: []int{2, 2, 2},
		},
		{
			I: []int{0, 1, 2},
			J: []int{0, 1, 2},
		},
		{
			I: []int{2, 1, 0},
			J: []int{0, 1, 2},
		},
	} {
		if x := checkSame(b, s.I, s.J); x != board.FillStateEmpty {
			return x
		}
	}

	return 0
}

func checkSame(b board.Board, i, j []int) board.FillState {
	init := b.Get(i[0], j[0])
	for x := 1; x < 3; x++ {
		if b.Get(i[x], j[x]) != init {
			return board.FillStateEmpty
		}
	}

	return init
}
