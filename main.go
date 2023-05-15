package main

import (
	"WP_projekt/client"
	"WP_projekt/game"
	"WP_projekt/menu"
	"WP_projekt/statystyki"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	play := func(isBot bool, targetNick string) {
		//g := &game.Game{Wpbot: true}
		g := &game.Game{Wpbot: isBot, TargetNick: targetNick}
		g.Start()
		g.WaitForBot()
		g.GetBoard()
		g.GetDescriptionsWithStatus()
		g.ShowGUI()

		status, _ := g.GetStatus()
		for strings.Compare(status.GameStatus, "game_in_progress") == 0 {
			if mojaTura := status.ShouldFire; !mojaTura {
				//fmt.Printf(status.GameStatus)
				fmt.Println("Czekam na ruch przeciwnika...")
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
	myMenu.AddOption("1", "Zagraj z botem", func() { play(true, "") })
	//myMenu.AddOption("1b", "Zagraj nie z botem", func() { play(false) })
	myMenu.AddOption("2", "Wyzwij przeciwnika", func() {
		err, players := client.List()
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		fmt.Printf("%v\n", players)

		if len(players) == 0 {
			fmt.Printf("Nie ma żadnego gracza czekającego na rozgrywkę :/\n")
			menu.WaitForClick()
		} else {
			defer time.Sleep(time.Second * 1)
			fmt.Printf("Lista graczy czekających na grę:\n")
			mapaNickow := make(map[string]string, 0)
			for i, v := range players {
				fmt.Printf("%v. %v\n", i, v)
				dupa := strconv.Itoa(i)
				mapaNickow[dupa] = v.Nick
			}
			fmt.Println()
			fmt.Printf("%#v\n", mapaNickow)
			fmt.Printf("Chcę zagrać z graczem nr:")
			reader := bufio.NewReader(os.Stdin)
			ans, _ := reader.ReadString('\n')
			ans = strings.Trim(ans, "\n")

			nick, ok := mapaNickow[ans]
			if !ok {
				fmt.Println("Źle wybrałeś gracza!")
				time.Sleep(time.Second * 2)
			} else {
				play(false, nick)
			}
		}
	})
	myMenu.AddOption("3", "Czekaj na wyzwanie", func() { play(false, "") })
	myMenu.AddOption("4", "Pokaż moja statystyki", func() { statystyki.GetInstance().PokazStroneSkutecznosci() })

	for {
		myMenu.DisplayMenu()
	}

}
