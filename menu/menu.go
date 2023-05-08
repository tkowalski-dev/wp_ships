package menu

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

/*
1. zagraj z botem
2. wyzwij przeciwnika
3. czekaj na wyzwanie
*/

type menu struct {
	//options map[string]func()
	options []menuOption
}

type menuOption struct {
	opt  string
	desc string
	fn   func()
}

func CreateMenu() *menu {
	m := &menu{}
	//m.options = make(map[string]func())
	//m.options["1"]
	m.options = make([]menuOption, 0)
	return m
}

func (m *menu) AddOption(option, desc string, fn func()) {
	//m.options[option] = fn
	m.options = append(m.options, menuOption{
		opt:  option,
		desc: desc,
		fn:   fn,
	})
}

func (m *menu) DisplayMenu() {
	header := ""
	reader := bufio.NewReader(os.Stdin)
MainLoop:
	for {
		fmt.Print("\033[H\033[2J")
		fmt.Println(header + "\nMenu:")
		for _, v := range m.options {
			fmt.Printf("%v. %v\n", v.opt, v.desc)
		}
		fmt.Printf("\nTwój wybór:")
		// czekaj na odpowiedź:
		ans, _ := reader.ReadString('\n')
		ans = strings.ReplaceAll(ans, "\n", "")
		for _, v := range m.options {
			if strings.Compare(v.opt, ans) == 0 {
				fmt.Print("\033[H\033[2J")
				v.fn()
				break MainLoop
			}
		}
		header = "Wybierz poprawną opcje menu!\n"
		//time.Sleep(time.Second * 2)
	}
}
