package handlers

import (
	"log"
	"testing"
)

func TestFetchDay(t *testing.T) {
	r, err := fetchDay()
	if err != nil {
		t.Error(err)
	}
	log.Println("fetchDay", r)
}
