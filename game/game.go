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
	SERVER_ADDR = "https://go-pjatk-server.fly.dev"
	INIT_GAME   = SERVER_ADDR + "/api/game"
	STATUS_GAME = SERVER_ADDR + "/api/game"
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
	Board          *gui.Board
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

func (g *Game) WaitForBot() bool {
	status, _ := g.GetStatus()
	for status.GameStatus == "waiting_wpbot" {
		time.Sleep(time.Second * 2)
		status, _ = g.GetStatus()
		//fmt.Printf("\n%#v\n", status)
	}
	return true
}

func (g *Game) PrepareGUI() {
	board := gui.New(
		gui.ConfigParams().
			HitChar('#').
			HitColor(color.FgRed).
			BorderColor(color.BgRed).
			RulerTextColor(color.BgYellow).NewConfig())
	g.Board = board
	//board.Import(g.lastStatusGame.OppShots)
	//board.Set(g.lastStatusGame.OppShots)

}
