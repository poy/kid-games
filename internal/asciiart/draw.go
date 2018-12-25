package asciiart

import (
	"fmt"
	"strings"
	"github.com/poy/kids-games/internal/board"
)

type Board interface {
	Get(i, j int) board.FillState
}

func Draw(b Board) string {
	var values []interface{}
	for j := 0; j < 3; j++ {
		for i := 0; i < 3; i++ {
			var s interface{}
			if v := b.Get(i, j); v != board.FillStateEmpty {
				s = v.String()
			} else {
				s = i + j*3 + 1
			}
			values = append(values, s)
		}
	}
	return fmt.Sprintf(strings.Replace(template, "s", " ", -1), values...)
}

var template = strings.TrimSpace(`
ss%vss|ss%vss|ss%vss
-----+-----+-----
ss%vss|ss%vss|ss%vss
-----+-----+-----
ss%vss|ss%vss|ss%vss
`)
