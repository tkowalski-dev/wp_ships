package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	SERVER_ADDR  = "https://go-pjatk-server.fly.dev"
	INIT_GAME    = SERVER_ADDR + "/api/game"
	STATUS_GAME  = SERVER_ADDR + "/api/game"
	BOARD        = SERVER_ADDR + "/api/game/board"
	DESCRIPTIONS = SERVER_ADDR + "/api/game/desc"
	FIRE         = SERVER_ADDR + "/api/game/fire"
	LIST         = SERVER_ADDR + "/api/game/list"
)

type client struct {
}

type List_Player struct {
	Game_status string `json:"game_status"`
	Nick        string `json:"nick"`
}

func List() (error, []List_Player) {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", LIST, nil)
	if err != nil {
		fmt.Printf("\n%v\n", err)
		return err, nil
	}

	//req.Header.Add("X-Auth-Token", g.AuthToken)
	req.Header.Add("Content-Type", "application/json")
	res, err := httpClient.Do(req)

	if err != nil {
		return err, nil
	}
	defer res.Body.Close()

	list_Players := make([]List_Player, 0)
	err = json.NewDecoder(res.Body).Decode(&list_Players)
	if err != nil {
		return err, nil
	}

	return nil, list_Players

}
