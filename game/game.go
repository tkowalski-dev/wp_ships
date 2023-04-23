package game

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	gui "github.com/grupawp/warships-lightgui"
	"net/http"
	"time"
)

const (
	SERVER_ADDR  = "https://go-pjatk-server.fly.dev"
	INIT_GAME    = SERVER_ADDR + "/api/game"
	STATUS_GAME  = SERVER_ADDR + "/api/game"
	BOARD        = SERVER_ADDR + "/api/game/board"
	DESCRIPTIONS = SERVER_ADDR + "/api/game/desc"
)

func CreateTestGame() *Game {
	return nil
}

type Game struct {
	Started        bool
	AuthToken      string
	Coords         []string `json:"coords"`
	Desc           string   `json:"desc"`
	Nick           string   `json:"nick"`
	TargetNick     string   `json:"target_nick"`
	Wpbot          bool     `json:"wpbot"`
	lastStatusGame StatusGame
	Board          []string `json:"board"`
	innerBoard     *gui.Board
	desc           string
	oppDesc        string
}

func (g *Game) GetLastStatusGame() StatusGame {
	return g.lastStatusGame
}

func (g *Game) Start() error {
	p := map[string]any{
		"wpbot": true,
	}

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

	return nil
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
	status, _ := g.GetStatus()
	for status.GameStatus == "waiting_wpbot" {
		time.Sleep(time.Second * 2)
		status, _ = g.GetStatus()
		//fmt.Printf("\n%#v\n", status)
	}
	return true
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
	board := gui.New(
		gui.ConfigParams().
			HitChar('#').
			HitColor(color.FgRed).
			BorderColor(color.BgRed).
			RulerTextColor(color.BgYellow).NewConfig())
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

func (g *Game) ShowGUI() {
	g.innerBoard.Display()
	fmt.Printf("Gracz: %16v %5v Przeciwnik: %11v", g.lastStatusGame.Nick, "", g.lastStatusGame.Opponent)
	//fmt.Printf("\n%v", g.desc)
	g.pokazOpisy()
}
