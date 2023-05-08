package main

import (
	"WP_projekt/game"
	"WP_projekt/menu"
	"WP_projekt/statystyki"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	play := func(isBot bool) {
		//g := &game.Game{Wpbot: true}
		g := &game.Game{Wpbot: isBot}
		g.Start()
		g.WaitForBot()
		g.GetBoard()
		g.GetDescriptionsWithStatus()
		g.ShowGUI()

		status, _ := g.GetStatus()
		for strings.Compare(status.GameStatus, "game_in_progress") == 0 {
			if mojaTura := status.ShouldFire; !mojaTura {
				fmt.Printf(status.GameStatus)
				time.Sleep(1 * time.Second)
			} else {
				//
				g.PobierzIWyswietlStrzalyPrzeciwnika()
				//g.ShowGUI()
				//g.PobierzStrzaly()
				g.WykonujStrzaly()
			}
			g.GetStatus()
		}
		fmt.Printf("'%v'", status.GameStatus)
		// Wyswietl informacje o zwyciezcy:
	}

	// menu:
	myMenu := menu.CreateMenu()
	myMenu.AddOption("0", "Zakoncz program", func() { os.Exit(0) })
	//myMenu.AddOption("1", "Zagraj z botem", playWithBot)
	myMenu.AddOption("1", "Zagraj z botem", func() { play(true) })
	//myMenu.AddOption("1b", "Zagraj nie z botem", func() { play(false) })
	myMenu.AddOption("2", "Wyzwij przeciwnika", func() {})
	myMenu.AddOption("3", "Czekaj na wyzwanie", func() { play(false) })
	myMenu.AddOption("4", "Pokaż moja statystyki", func() { statystyki.GetInstance().PokazStroneSkutecznosci() })

	for {
		myMenu.DisplayMenu()
	}

}
