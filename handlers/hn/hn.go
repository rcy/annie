package hn

import (
	"encoding/json"
	"fmt"
	"goirc/bot"
	"math/rand"
	"net/http"
)

func Handle(params bot.HandlerParams) error {
	resp, err := http.Get("https://hacker-news.firebaseio.com/v0/beststories.json")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	ids := []int{}

	err = json.NewDecoder(resp.Body).Decode(&ids)
	if err != nil {
		return err
	}

	id := ids[rand.Intn(30)]

	resp, err = http.Get(fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", id))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var item struct {
		Title string `json:"title"`
		URL   string `json:"url"`
	}
	err = json.NewDecoder(resp.Body).Decode(&item)
	if err != nil {
		return err
	}
	params.Privmsgf(params.Target, "%s", item.Title)

	return nil
}
