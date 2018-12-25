package game

import (
	"errors"
	"sync"
	"github.com/poy/kids-games/internal/board"
)

type Game struct {
	mu   sync.RWMutex
	turn int
	b    board.Board
}

func New(b board.Board) *Game {
	return &Game{
		b: b,
	}
}

func (g *Game) P1(x int) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.turn%2 != 0 {
		return errors.New("out of turn")
	}

	i, j := g.coordinates(x)
	if !g.b.Set(i, j, board.FillStateX) {
		return errors.New("spot taken")
	}
	g.turn++

	return nil
}

func (g *Game) P2(x int) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.turn%2 != 1 {
		return errors.New("out of turn")
	}

	i, j := g.coordinates(x)
	if !g.b.Set(i, j, board.FillStateO) {
		return errors.New("spot taken")
	}
	g.turn++
	return nil
}

func (g *Game) State() (turn int, winner board.FillState) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.turn, Winner(g.b)
}

func (g *Game) Board() board.BoardReader {
	return g.b
}

func (g *Game) coordinates(x int) (i, j int) {
	x--
	return x % 3, x / 3
}
