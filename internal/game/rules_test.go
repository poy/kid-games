package game_test

import (
	"github.com/poy/kids-games/internal/asciiart"
	"github.com/poy/kids-games/internal/board"
	"github.com/poy/kids-games/internal/game"
	"testing"
)

func TestRules(t *testing.T) {
	for _, test := range []struct {
		Name   string
		Board  board.Board
		Winner board.FillState
	}{
		{
			Name: "empty",
		},
		{
			Name: "cat",
			Board: newBB().
				O().X().O().
				O().X().X().
				X().O().X().b,
		},
		{
			Name:   "O Left Down",
			Winner: board.FillStateO,
			Board: newBB().
				O().X().E().
				O().X().X().
				O().O().X().b,
		},
		{
			Name:   "X Middle Down",
			Winner: board.FillStateX,
			Board: newBB().
				O().X().X().
				X().X().O().
				O().X().X().b,
		},
		{
			Name:   "X Right Down",
			Winner: board.FillStateX,
			Board: newBB().
				O().X().X().
				X().O().X().
				O().X().X().b,
		},
		{
			Name:   "O Upper Across",
			Winner: board.FillStateO,
			Board: newBB().
				O().O().O().
				X().X().O().
				O().X().X().b,
		},
		{
			Name:   "X Middle Across",
			Winner: board.FillStateX,
			Board: newBB().
				X().O().O().
				X().X().X().
				O().X().O().b,
		},
		{
			Name:   "O Lower Across",
			Winner: board.FillStateO,
			Board: newBB().
				X().O().O().
				O().X().X().
				O().O().O().b,
		},
		{
			Name:   "O Diag L->R",
			Winner: board.FillStateO,
			Board: newBB().
				O().X().O().
				X().O().X().
				X().X().O().b,
		},
		{
			Name:   "X Diag R->L",
			Winner: board.FillStateX,
			Board: newBB().
				O().X().X().
				X().X().O().
				X().O().O().b,
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			if test.Board == nil {
				test.Board = board.New()
			}

			actual := game.Winner(test.Board)
			if actual != test.Winner {
				t.Fatalf("expected %q to equal %q:\n\n%s", actual, test.Winner, asciiart.Draw(test.Board))
			}
		})
	}
}

type boardBuilder struct {
	b    board.Board
	i, j int
}

func newBB() *boardBuilder {
	return &boardBuilder{
		b: board.New(),
	}
}

func (b *boardBuilder) X() *boardBuilder {
	b.b.Set(b.i, b.j, board.FillStateX)
	b.adv()
	return b
}

func (b *boardBuilder) O() *boardBuilder {
	b.b.Set(b.i, b.j, board.FillStateO)
	b.adv()
	return b
}

func (b *boardBuilder) E() *boardBuilder {
	b.b.Set(b.i, b.j, board.FillStateEmpty)
	b.adv()
	return b
}

func (b *boardBuilder) adv() {
	b.i++
	if b.i == 3 {
		b.i = 0
		b.j++
	}
}
