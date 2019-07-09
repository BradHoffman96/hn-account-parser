package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

//Account struct for HN Users
type Account struct {
	About     string `json:"about"`
	Created   int    `json:"created"`
	ID        string `json:"id"`
	Karma     int    `json:"karma"`
	Submitted []int  `json:"submitted"`
}

var client = &http.Client{Timeout: 10 * time.Second}

func main() {
	account := Account{}
	getAccount("jandrewrogers", &account)
	fmt.Println(account.Submitted)
}

func getAccount(id string, target interface{}) error {
	url := []string{"https://hacker-news.firebaseio.com/v0/user/", id, ".json?print=pretty"}

	resp, err := client.Get(strings.Join(url, ""))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}
