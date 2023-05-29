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
	nick := ""

	play := func(isBot bool, targetNick string) {
		//g := &game.Game{Wpbot: true}
		g := &game.Game{Wpbot: isBot, TargetNick: targetNick, MyNick: nick}
		g.Start()
		g.WaitForBot()
		g.GetBoard()
		g.GetDescriptionsWithStatus()
		g.ShowGUI()

		status, _ := g.GetStatus()
		for status.GameStatus == "game_in_progress" {
			if mojaTura := status.ShouldFire; !mojaTura {
				fmt.Printf(status.GameStatus)
				fmt.Printf("%+v\n", status)
				fmt.Println("Czekam na ruch przeciwnika...")
				time.Sleep(1 * time.Second)
			} else {
				//
				g.PobierzIWyswietlStrzalyPrzeciwnika()
				//g.ShowGUI()
				//g.PobierzStrzaly()
				g.WykonujStrzaly(nil)
			}
			status, _ = g.GetStatus()
			//fmt.Printf("%v\n", g.GetLastStatusGame())
			//time.Sleep(time.Second * 5)
		}
		fmt.Printf("'%v'\n", status.GameStatus)
		fmt.Printf("%v\n\n", status)
		fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		// Wyswietl informacje o zwyciezcy:
		fmt.Printf("Wynik rozgrywki:%v\n", status.GameStatus)
		time.Sleep(time.Second * 3)
		menu.WaitForClick()
	}

	playAutomatic := func(isBot bool, targetNick string) {
		idMove := 0
		moves := make([]string, 0, 16)
		moves = append(moves, "A2", "A4", "A6", "A8", "A10")
		moves = append(moves, "B1", "B3", "B5", "B7", "B9")
		moves = append(moves, "C2", "C4", "C6", "C8", "C10")
		moves = append(moves, "D1", "D3", "D5", "D7", "D9")
		moves = append(moves, "E2", "E4", "E6", "E8", "E10")
		moves = append(moves, "F1", "F3", "F5", "F7", "F9")
		moves = append(moves, "G2", "G4", "G6", "G8", "G10")
		moves = append(moves, "H1", "H3", "H5", "H7", "H9")
		moves = append(moves, "I2", "I4", "I6", "I8", "I10")
		//moves = append(moves, "J1", "J3", "J5", "J7", "J9")

		//g := &game.Game{Wpbot: true}
		g := &game.Game{Wpbot: isBot, TargetNick: targetNick, MyNick: nick}
		g.Start()
		g.WaitForBot()
		g.GetBoard()
		g.GetDescriptionsWithStatus()
		g.ShowGUI()

		status, _ := g.GetStatus()
		for status.GameStatus == "game_in_progress" {
			if mojaTura := status.ShouldFire; !mojaTura {
				fmt.Printf(status.GameStatus)
				fmt.Printf("%+v\n", status)
				fmt.Println("Czekam na ruch przeciwnika...")
				time.Sleep(1 * time.Second)
			} else {
				//
				if idMove < len(moves) {
					g.PobierzIWyswietlStrzalyPrzeciwnika()
					var strzal string = moves[idMove]
					idMove++
					g.WykonujStrzaly(&strzal)
					time.Sleep(1 * time.Second)
					//status, _ = g.GetStatus()
					g.ShowGUI()
				} else {
					g.PobierzIWyswietlStrzalyPrzeciwnika()
					//g.ShowGUI()
					//g.PobierzStrzaly()
					g.WykonujStrzaly(nil)
				}
			}
			status, _ = g.GetStatus()
			//fmt.Printf("%v\n", g.GetLastStatusGame())
			//time.Sleep(time.Second * 5)
		}
		fmt.Printf("'%v'\n", status.GameStatus)
		fmt.Printf("%v\n\n", status)
		fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		// Wyswietl informacje o zwyciezcy:
		fmt.Printf("Wynik rozgrywki:%v\n", status.GameStatus)
		time.Sleep(time.Second * 3)
		menu.WaitForClick()
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
				lala := strconv.Itoa(i)
				mapaNickow[lala] = v.Nick
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
				time.Sleep(time.Second * 1)
			} else {
				play(false, nick)
			}
		}
	})
	myMenu.AddOption("3", "Czekaj na wyzwanie", func() { play(false, "") })
	myMenu.AddOption("4", "Pokaż moja statystyki", func() { statystyki.GetInstance().PokazStroneSkutecznosci() })
	myMenu.AddOption("5", "Zagraj z botem automatycznie", func() { playAutomatic(true, "") })

	// Nick req:
	ustawNick := func() {
		fmt.Print("\033[H\033[2J")
		fmt.Printf("\nUstaw mój nick na: ")
		reader := bufio.NewReader(os.Stdin)
		ans, _ := reader.ReadString('\n')
		ans = strings.Trim(ans, "\n")
		if len(ans) == 0 {
			fmt.Println("Twój nick zostanie przydzielony automatycznie.")
			nick = ""
			fmt.Println("Nick będzie ustawiany automatycznie.")
		} else {
			nick = ans
			fmt.Println("Nick został ustawiony!")
		}
		time.Sleep(time.Second * 2)
	}
	myMenu.AddOption("6", "Ustaw nick", func() { ustawNick() })
	ustawNick()

	for {
		myMenu.DisplayMenu(nick)
	}

}
