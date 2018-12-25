package main

import (
	"fmt"
	"github.com/poy/kids-games/internal/asciiart"
	"github.com/poy/kids-games/internal/board"
	"github.com/poy/kids-games/internal/game"
)

func main() {
	g := game.New(board.New())

	for i := 0; i < 9; i++ {
		play := g.P1
		if i%2 != 0 {
			play = g.P2
		}

		fmt.Printf("\n%s\n\n", asciiart.Draw(g.Board()))
		var entry uint
		for {
			fmt.Printf("Player %d: ", i%2+1)

			_, err := fmt.Scanln(&entry)
			if err != nil || entry == 0 || entry > 9 {
				fmt.Println("invalid... choose a number from the board")
				continue
			}

			if err := play(int(entry)); err != nil {
				fmt.Println(err)
				continue
			}

			if _, winner := g.State(); winner != board.FillStateEmpty {
				fmt.Printf("\n%s\n%v's wins!\n", asciiart.Draw(g.Board()), winner)
				return
			}

			break
		}
	}
}
