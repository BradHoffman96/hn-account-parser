package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

//Comment struct for HN
type Comment struct {
	ID     int    `json:"id"`
	By     string `json:"by"`
	Type   string `json:"type"`
	Time   int    `json:"time"`
	Text   string `json:"text"`
	Parent int    `json:"parent"`
	Kids   []int  `json:"kids"`
	URL    string `json:"url"`
}

var comments []Comment
var client = &http.Client{Timeout: 10 * time.Second}

func main() {
	account := Account{}
	getAccount("jandrewrogers", &account)
	getComments(account)

	fmt.Println(len(comments))
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

func getComments(account Account) {
	for _, item := range account.Submitted {
		url := []string{"https://hacker-news.firebaseio.com/v0/item/", strconv.Itoa(item), ".json?print=pretty"}

		resp, err := client.Get(strings.Join(url, ""))
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		comment := Comment{}
		err = json.NewDecoder(resp.Body).Decode(&comment)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(comment.Text)

		comments = append(comments, comment)
	}
}
