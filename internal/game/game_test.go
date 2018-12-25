package game_test

import (
	"github.com/poy/kids-games/internal/board"
	"github.com/poy/kids-games/internal/game"
	"testing"
)

// 1 | 2 | 3
// --+---+---
// 4 | 5 | 6
// --+---+---
// 7 | 8 | 9
func TestGame(t *testing.T) {
	for _, test := range []struct {
		Name    string
		Errored bool
		Game    *gamePlayer

		Turn   int
		Winner board.FillState
	}{
		{
			Name:   "player 1 wins",
			Winner: board.FillStateX,
			Game:   newGP().P1(1).P2(2).P1(4).P2(5).P1(7),
			Turn:   5,
		},
		{
			Name:   "player 2 wins",
			Winner: board.FillStateO,
			Game:   newGP().P1(1).P2(5).P1(9).P2(7).P1(4).P2(3),
			Turn:   6,
		},
		{
			Name:    "player 2 went first",
			Errored: true,
			Game:    newGP().P2(1),
			Turn:    0,
		},
		{
			Name:    "player 1 went twice",
			Errored: true,
			Game:    newGP().P1(1).P1(2),
			Turn:    1,
		},
		{
			Name:    "player 2 went twice",
			Errored: true,
			Game:    newGP().P1(9).P2(1).P2(2),
			Turn:    2,
		},
		{
			Name:    "player 1 spot is taken",
			Errored: true,
			Game:    newGP().P1(9).P2(1).P1(1),
			Turn:    2,
		},
		{
			Name:    "player 2 spot is taken",
			Errored: true,
			Game:    newGP().P1(1).P2(1),
			Turn:    1,
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			for _, act := range test.Game.actions {
				act()
			}
			if test.Errored && test.Game.err == nil {
				t.Fatal("expected an error")
			}

			if !test.Errored && test.Game.err != nil {
				t.Fatal(test.Game.err)
			}

			turn, winner := test.Game.g.State()
			if test.Turn != turn {
				t.Fatalf("expected turn %d to equal %d", turn, test.Turn)
			}

			if test.Winner != winner {
				t.Fatalf("expected winner %q to equal %q", winner, test.Winner)
			}
		})
	}
}

func TestGameRace(t *testing.T) {
	g := game.New(board.New())

	go func() {
		for i := 0; i < 10; i++ {
			g.P1(5)
		}
	}()

	go func() {
		for i := 0; i < 10; i++ {
			g.P2(5)
		}
	}()

	go func() {
		for i := 0; i < 10; i++ {
			g.State()
		}
	}()

	g.P1(5)
	g.P1(5)
}

type gamePlayer struct {
	g       *game.Game
	actions []func()
	err     error
}

func newGP() *gamePlayer {
	return &gamePlayer{
		g: game.New(board.New()),
	}
}

func (p *gamePlayer) P1(x int) *gamePlayer {
	p.actions = append(p.actions, func() {
		if p.err != nil {
			return
		}
		p.err = p.g.P1(x)
	})
	return p
}

func (p *gamePlayer) P2(x int) *gamePlayer {
	p.actions = append(p.actions, func() {
		if p.err != nil {
			return
		}
		p.err = p.g.P2(x)
	})
	return p
}
