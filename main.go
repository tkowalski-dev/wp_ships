package main

import (
	"WP_projekt/game"
)

func main() {
	g := &game.Game{Wpbot: true}
	g.Start()
	g.WaitForBot()
	g.PrepareGUI()
	g.Board.Display()

}
