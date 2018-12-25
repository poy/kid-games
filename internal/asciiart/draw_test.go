package asciiart_test

import (
	"strings"
	"github.com/poy/kids-games/internal/asciiart"
	"github.com/poy/kids-games/internal/board"
	"testing"
)

func TestDraw(t *testing.T) {
	for _, test := range []struct {
		Name     string
		Expected string
		Board    stubBoard
	}{
		{
			Name: "empty",
			Expected: strings.TrimSpace(`
  1  |  2  |  3
-----+-----+-----
  4  |  5  |  6
-----+-----+-----
  7  |  8  |  9
`),
		},
		{
			Name: "corners X",
			Board: stubBoard{
				pair{I: 0, J: 0}: board.FillStateX,
				pair{I: 2, J: 0}: board.FillStateX,
				pair{I: 0, J: 2}: board.FillStateX,
				pair{I: 2, J: 2}: board.FillStateX,
			},
			Expected: strings.TrimSpace(`
  X  |  2  |  X
-----+-----+-----
  4  |  5  |  6
-----+-----+-----
  X  |  8  |  X
`),
		},
		{
			Name: "middle O",
			Board: stubBoard{
				pair{I: 1, J: 0}: board.FillStateO,
				pair{I: 0, J: 1}: board.FillStateO,
				pair{I: 1, J: 1}: board.FillStateO,
				pair{I: 2, J: 1}: board.FillStateO,
				pair{I: 1, J: 2}: board.FillStateO,
			},
			Expected: strings.TrimSpace(`
  1  |  O  |  3
-----+-----+-----
  O  |  O  |  O
-----+-----+-----
  7  |  O  |  9
`),
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			actual := asciiart.Draw(test.Board)
			if remSpace(actual) != remSpace(test.Expected) {
				t.Fatalf("expected:\n%s\nto equal:\n%s", actual, test.Expected)
			}
			if len(actual) != 89 {
				t.Fatalf("wrong spacing: %d\n%s", len(actual), actual)
			}
			t.Logf("result:\n%s", actual)
		})
	}
}

func remSpace(s string) string {
	return strings.Replace(s, " ", "", -1)
}

type pair struct {
	I int
	J int
}

type stubBoard map[pair]board.FillState

func (b stubBoard) Get(i, j int) board.FillState {
	for k, v := range b {
		if k.I == i && k.J == j {
			return v
		}
	}

	return board.FillStateEmpty
}
