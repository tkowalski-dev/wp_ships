package game

import (
	"WP_projekt/statystyki"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	gui "github.com/grupawp/warships-lightgui/v2"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	SERVER_ADDR  = "https://go-pjatk-server.fly.dev"
	INIT_GAME    = SERVER_ADDR + "/api/game"
	STATUS_GAME  = SERVER_ADDR + "/api/game"
	BOARD        = SERVER_ADDR + "/api/game/board"
	DESCRIPTIONS = SERVER_ADDR + "/api/game/desc"
	FIRE         = SERVER_ADDR + "/api/game/fire"
	REFRESH_GAME = SERVER_ADDR + "/api/game/refresh"
)

func CreateTestGame() *Game {
	return nil
}

type Game struct {
	Started            bool
	AuthToken          string
	Coords             []string `json:"coords"`
	Desc               string   `json:"desc"`
	Nick               string   `json:"nick"`
	TargetNick         string   `json:"target_nick"`
	Wpbot              bool     `json:"wpbot"`
	lastStatusGame     StatusGame
	Board              []string `json:"board"`
	innerBoard         *gui.Board
	desc               string
	oppDesc            string
	strzalyPrzeciwnika []string
	stats              *statystyki.Statystyki
}

func (g *Game) GetLastStatusGame() StatusGame {
	return g.lastStatusGame
}

func (g *Game) Start() error {
	p := map[string]any{
		"wpbot":       g.Wpbot,
		"target_nick": g.TargetNick,
	}
	g.stats = statystyki.GetInstance()

	buff := &bytes.Buffer{}
	err := json.NewEncoder(buff).Encode(p)
	if err != nil {
		return err
	}

	hClient := http.Client{
		Timeout: 10 * time.Second,
	}

	response, err := hClient.Post(INIT_GAME, "application/json", buff)
	if err != nil {
		return err
	}

	fmt.Printf("%#v\n", response)

	token := response.Header.Get("X-Auth-Token")
	fmt.Printf("\n%#v\n", token)

	if len(token) == 0 {
		return fmt.Errorf("Nie otrzymano tokena!")
	}

	g.AuthToken = token
	g.Started = true

	g.strzalyPrzeciwnika = make([]string, 0)
	g.PrepareGUI()

	return nil
}

func (g *Game) PobierzIWyswietlStrzalyPrzeciwnika() {
	for i := len(g.strzalyPrzeciwnika); i < len(g.lastStatusGame.OppShots); i++ {
		//TODO sprawdzić trafienie:
		g.innerBoard.Set(gui.Left, g.lastStatusGame.OppShots[i], gui.Hit)
		g.innerBoard.Set(gui.Left, g.lastStatusGame.OppShots[i], gui.Miss)
	}
	g.strzalyPrzeciwnika = g.lastStatusGame.OppShots
}

func (g *Game) getHTTP() {

}

type StatusGame struct {
	Desc           string   `json:"desc"`
	GameStatus     string   `json:"game_status"`
	LastGameStatus string   `json:"last_game_status"`
	Nick           string   `json:"nick"`
	OppDesc        string   `json:"opp_desc"`
	OppShots       []string `json:"opp_shots"`
	Opponent       string   `json:"opponent"`
	ShouldFire     bool     `json:"should_fire"`
	Timer          int      `json:"timer"`
}

func (g *Game) GetStatus() (StatusGame, error) {
	sg := &StatusGame{}

	httpClient := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", STATUS_GAME, nil)
	if err != nil {
		fmt.Printf("\n%v\n", err)
		return *sg, err
	}

	req.Header.Add("X-Auth-Token", g.AuthToken)
	req.Header.Add("Content-Type", "application/json")
	res, err := httpClient.Do(req)

	if err != nil {
		fmt.Printf("\n%v\n", err)
		return *sg, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(sg)
	if err != nil {
		fmt.Printf("\n%v\n", err)
		return *sg, err
	}

	g.lastStatusGame = *sg
	return *sg, err
}

func (g *Game) RefreshGame() error {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", REFRESH_GAME, nil)
	if err != nil {
		fmt.Printf("\n%v\n", err)
		return err
	}

	req.Header.Add("X-Auth-Token", g.AuthToken)
	req.Header.Add("Content-Type", "application/json")
	res, err := httpClient.Do(req)

	if err != nil {
		fmt.Printf("\n%v\n", err)
		return err
	}
	defer res.Body.Close()

	return err
}

func (g *Game) GetBoard() (any, error) {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", BOARD, nil)
	if err != nil {
		fmt.Printf("\n%v\n", err)
		return nil, err
	}

	req.Header.Add("X-Auth-Token", g.AuthToken)
	req.Header.Add("Content-Type", "application/json")
	res, err := httpClient.Do(req)

	if err != nil {
		fmt.Printf("\n%v\n", err)
		return nil, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&g)
	if err != nil {
		fmt.Printf("\n%v\n", err)
		return nil, err
	}

	for _, v := range g.Board {
		g.innerBoard.Set(gui.Left, v, gui.Ship)
	}
	return g.Board, err
}

func (g *Game) WaitForBot() bool {
	tickerRefresh := time.NewTicker(time.Second * 10)
	tickerMessage := time.After(time.Duration(0))

	for {
		select {
		case <-tickerMessage:
			tickerMessage = time.After(time.Second * 2)
			status, _ := g.GetStatus()
			switch status.GameStatus {
			case "waiting_wpbot":
				fmt.Printf("Czekam na wpbota...\n")
			case "waiting":
				fmt.Printf("Czekam na wyzwanie...\n")
			default:
				return true
			}
		case <-tickerRefresh.C:
			_ = g.RefreshGame()
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}

	////for status.GameStatus == "waiting_wpbot" {
	//for g.Wpbot && status.GameStatus == "waiting_wpbot" {
	//	fmt.Printf("Czekam na wpbota...\n")
	//	time.Sleep(time.Second * 2)
	//	status, _ = g.GetStatus()
	//	//fmt.Printf("\n%#v\n", status)
	//}
	//for !g.Wpbot && status.GameStatus == "waiting" {
	//	fmt.Printf("Czekam na wyzwanie...\n")
	//	time.Sleep(time.Second * 2)
	//	status, _ = g.GetStatus()
	//	//fmt.Printf("\n%#v\n", status)
	//}
	//return true
}

func (g *Game) GetDescriptionsWithStatus() (StatusGame, error) {
	sg := &StatusGame{}

	httpClient := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", DESCRIPTIONS, nil)
	if err != nil {
		fmt.Printf("\n%v\n", err)
		return *sg, err
	}

	req.Header.Add("X-Auth-Token", g.AuthToken)
	req.Header.Add("Content-Type", "application/json")
	res, err := httpClient.Do(req)

	if err != nil {
		fmt.Printf("\n%v\n", err)
		return *sg, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(sg)
	if err != nil {
		fmt.Printf("\n%v\n", err)
		return *sg, err
	}

	g.lastStatusGame = *sg
	g.desc = g.lastStatusGame.Desc
	g.oppDesc = g.lastStatusGame.OppDesc
	return *sg, err
}

func (g *Game) PrepareGUI() {
	cfg := gui.NewConfig()
	cfg.HitChar = '#'
	cfg.HitColor = color.FgRed
	cfg.BorderColor = color.BgRed
	cfg.RulerTextColor = color.BgYellow
	board := gui.New(cfg)
	board.Display()

	g.innerBoard = board
}

func (g *Game) pokazOpisy() {
	moj := make([]string, 0)
	wrog := make([]string, 0)
	for i := 0; i < len([]rune(g.desc)); i += 27 {
		od := i
		do := i + 27
		if do >= len([]rune(g.desc)) {
			do = len([]rune(g.desc)) - 1
		}
		moj = append(moj, string([]rune(g.desc)[od:do]))
	}
	for i := 0; i < len([]rune(g.oppDesc)); i += 27 {
		od := i
		do := i + 27
		if do >= len([]rune(g.oppDesc)) {
			do = len([]rune(g.oppDesc)) - 1
		}
		wrog = append(wrog, string([]rune(g.oppDesc)[od:do]))
	}
	//fmt.Printf("%#v", moj)
	//fmt.Printf("\n%#v", wrog)
	deli := "-"
	linelength := 5 + 2*27
	line := ""
	for i := 0; i < linelength; i++ {
		line += deli
	}
	deli = "|"

	fmt.Printf("\n%v", line)
	for i := 0; i < len(moj) || i < len(wrog); i++ {
		l, r := "", ""
		if i < len(moj) {
			l = moj[i]
		}
		if i < len(wrog) {
			r = wrog[i]
		}
		fmt.Printf("\n%v%27v%v %v%27v%v", deli, l, deli, deli, r, deli)
	}
	fmt.Printf("\n%v", line)
}

func (g *Game) PobierzStrzaly() string {
	//g.ShowGUI()
	czyPoprawne := false
	czyBlad := false
	coords := ""
	for !czyPoprawne {
		g.ShowGUI()
		if czyBlad {
			fmt.Printf("\nNiepoprawne dane, podaj współrzędne A1-J10!")
		}
		fmt.Printf("\nPodaj współrzędne strzału:")
		czyBlad = true
		reader := bufio.NewReader(os.Stdin)
		ans, _ := reader.ReadString('\n')
		//r, _ := regexp.Compile("[a-jA-J][]")
		//fmt.Printf("odp: %v", ans)
		if len(ans) != 3 && len(ans) != 4 {
			continue
		}
		//if !(strings.HasPrefix(ans, "[a-z]") && int(([]rune(ans))[0:1]) <= int('J')) {
		firstChar := []rune(ans)[0]
		next := []rune(ans)[1 : len([]rune(ans))-1]
		//fmt.Println(firstChar, " ", next)

		if !((firstChar >= 'A' && firstChar <= 'J') || (firstChar >= 'a' && firstChar <= 'j')) {
			continue
		}
		y, err := strconv.ParseInt(string(next), 10, 32)
		if err != nil {
			continue
		}
		if y < 1 || y > 10 {
			continue
		}
		//time.Sleep(2 * time.Second)
		czyPoprawne = true
		coords = string(firstChar) + strconv.Itoa(int(y))
	}
	fmt.Printf("\nStrzelam...")
	return coords
}

func (g *Game) WykonujStrzaly(pre *string) {
	var strzal string
	if pre == nil {
		strzal = g.PobierzStrzaly()
	} else {
		strzal = *pre
	}

	p := map[string]any{
		"coord": strzal,
	}

	buff := &bytes.Buffer{}
	err := json.NewEncoder(buff).Encode(p)

	httpClient := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", FIRE, buff)
	if err != nil {
		fmt.Printf("\n%v\n", err)
		return
	}

	req.Header.Add("X-Auth-Token", g.AuthToken)
	req.Header.Add("Content-Type", "application/json")
	//res, err := httpClient.Do(req)
	respr, _ := httpClient.Do(req)
	defer respr.Body.Close()
	reader := bufio.NewReader(respr.Body)
	//resp, _ := io.ReadAll(reader)
	//fmt.Println(string(resp))
	result := struct {
		Result string `json:"result"`
	}{}
	json.NewDecoder(reader).Decode(&result)
	fmt.Printf("\n%+v", result)
	time.Sleep(time.Second * 1)
	//if strings.Compare(result.Result, "hit") == 0 {
	//	g.innerBoard.Set(gui.Right, strzal, gui.Hit)
	//} else if strings.Compare(result.Result, "miss") == 0 {
	//	g.innerBoard.Set(gui.Right, strzal, gui.Miss)
	//} else {
	//	//sunk
	//	g.innerBoard.Set(gui.Right, strzal, gui.Hit)
	//	g.innerBoard.CreateBorder(gui.Right, strzal)
	//}
	switch result.Result {
	case "hit":
		g.innerBoard.Set(gui.Right, strzal, gui.Hit)
	case "miss":
		g.innerBoard.Set(gui.Right, strzal, gui.Miss)
	case "sunk":
		g.innerBoard.Set(gui.Right, strzal, gui.Hit)
		g.innerBoard.CreateBorder(gui.Right, strzal)
	default:
		return
	}
	//time.Sleep(5 * time.Second)
}

func (g *Game) ShowGUI() {
	g.innerBoard.Display()
	fmt.Printf("Gracz: %16v %5v Przeciwnik: %11v", g.lastStatusGame.Nick, "", g.lastStatusGame.Opponent)
	//fmt.Printf("\n%v", g.desc)
	g.pokazOpisy()
	if g.lastStatusGame.ShouldFire {
		fmt.Printf("\nTwój ruch!\n")
	}
}
