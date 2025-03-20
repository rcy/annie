package summary

import (
	"context"
	"fmt"
	"goirc/db/model"
	db "goirc/model"
	"testing"
	"time"
)

func TestHTML(t *testing.T) {
	q := model.New(db.DB)
	s := New(q, time.Now().Add(-time.Hour*24*7), time.Now().Add(-time.Hour*24*0))
	// s.SetCompleteFn(func(context.Context, string, string, string) (string, error) {
	// 	return "mock ai completed text", nil
	// })
	// s.SetGetTitleFn(func(string) (string, error) {
	// 	return "mock title", nil
	// })

	err := s.LoadAll(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	bytes, err := s.HTML(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%s\n", string(bytes))
}
