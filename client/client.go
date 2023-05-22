package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

type Client struct {
	timeout *time.Duration
}

// func (c *Client) GET(url string, token *string) (*[]byte, error) {
func (c *Client) GET(url string, token *string) (*string, error) {
	localTimeout := 10 * time.Second
	if c.timeout != nil {
		localTimeout = *c.timeout
	}
	httpClient := &http.Client{Timeout: localTimeout}
	//req, err := http.NewRequest("GET", LIST, nil)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if token != nil {
		req.Header.Add("X-Auth-Token", *token)
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := httpClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	body := string(bytes)

	return &body, nil
	//return &bytes, nil
}

// func (c *Client) POST(url string, bodyToSend interface{}, token *string) (*[]byte, error) {
func (c *Client) POST(url string, bodyToSend interface{}, token *string) (*string, error) {
	//p := map[string]any{
	//	"coord": strzal,
	//}

	buff := &bytes.Buffer{}
	err := json.NewEncoder(buff).Encode(bodyToSend)
	if err != nil {
		return nil, err
	}

	localTimeout := 10 * time.Second
	if c.timeout != nil {
		localTimeout = *c.timeout
	}
	httpClient := &http.Client{Timeout: localTimeout}
	//req, err := http.NewRequest("POST", FIRE, buff)
	req, err := http.NewRequest("POST", url, buff)
	if err != nil {
		return nil, err
	}

	if token != nil {
		req.Header.Add("X-Auth-Token", *token)
	}
	req.Header.Add("Content-Type", "application/json")
	respr, _ := httpClient.Do(req)
	defer respr.Body.Close()

	bajty, err := io.ReadAll(respr.Body)
	if err != nil {
		return nil, err
	}
	body := string(bajty)

	return &body, nil
	//return &bajty, nil
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
