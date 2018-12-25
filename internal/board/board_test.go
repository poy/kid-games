package board_test

import (
	"github.com/poy/kids-games/internal/board"
	"testing"
)

func TestBoard(t *testing.T) {
	for _, test := range []struct {
		Name string
		I    int
		J    int
		Fill board.FillState
	}{
		{
			Name: "upper left corner",
			I:    0,
			J:    0,
			Fill: board.FillStateX,
		},
		{
			Name: "upper right corner",
			I:    2,
			J:    0,
			Fill: board.FillStateX,
		},
		{
			Name: "bottom left corner",
			I:    0,
			J:    2,
			Fill: board.FillStateX,
		},
		{
			Name: "bottom right corner",
			I:    2,
			J:    2,
			Fill: board.FillStateX,
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			b := board.New()
			b.Set(test.I, test.J, test.Fill)
			if b.Get(test.I, test.J) != test.Fill {
				t.Fatal("wrong")
			}
		})
	}
}

func TestBoardAlreadyFilled(t *testing.T) {
	b := board.New()
	if !b.Set(0, 0, board.FillStateO) {
		t.Fatal("expected to be true")
	}
	if b.Set(0, 0, board.FillStateO) {
		t.Fatal("expected to be false")
	}
}

func TestBoardRaceDetector(t *testing.T) {
	b := board.New()

	go func() {
		for i := 0; i < 10; i++ {
			b.Set(0, 0, board.FillStateO)
		}
	}()

	go func() {
		for i := 0; i < 10; i++ {
			b.Get(0, 0)
		}
	}()

	b.Set(0, 0, board.FillStateO)
	b.Get(0, 0)
}

func TestBoardPanics(t *testing.T) {
	for _, test := range []struct {
		Name string
		I    int
		J    int
	}{
		{
			Name: "negative i",
			I:    -1,
			J:    0,
		},
		{
			Name: "negative j",
			I:    0,
			J:    -1,
		},
		{
			Name: "large i",
			I:    3,
			J:    0,
		},
		{
			Name: "large j",
			I:    0,
			J:    3,
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			b := board.New()

			f := func(f func()) interface{} {
				var err interface{}
				func() {
					defer func() {
						err = recover()
					}()

					f()
				}()

				return err
			}

			if err := f(func() { b.Set(test.I, test.J, board.FillStateO) }); err == nil {
				t.Fatal("expected to panic")
			}

			if err := f(func() { b.Get(test.I, test.J) }); err == nil {
				t.Fatal("expected to panic")
			}
		})
	}
}
