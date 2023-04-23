package main

import (
	"WP_projekt/game"
)

func main() {
	g := &game.Game{Wpbot: true}
	g.Start()
	g.PrepareGUI()
	g.WaitForBot()
	g.GetBoard()
	g.GetDescriptionsWithStatus()
	g.GetStatus()
	g.ShowGUI()

}
