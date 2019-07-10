package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
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
	ID     int             `json:"id"`
	By     string          `json:"by"`
	Type   string          `json:"type"`
	Time   int             `json:"time"`
	Text   json.RawMessage `json:"text"`
	Parent int             `json:"parent"`
	Kids   []int           `json:"kids"`
	URL    string          `json:"url"`
}

var comments []Comment
var client = &http.Client{Timeout: 10 * time.Second}

func main() {
	account := Account{}

	getAccount("jandrewrogers", &account)
	getComments(account)
	writeToFile()
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

		if strings.Contains(string(comment.Text), "spatial") {
			comments = append(comments, comment)
		}
	}
}

func writeToFile() {
	commentsJSON, err := JSONMarshal(comments)
	if err != nil {
		log.Fatal(err)
	}

	commentsJSON = []byte(strings.Replace(string(commentsJSON), "&#x27;", "'", -1))

	err = ioutil.WriteFile("output.json", commentsJSON, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}
