package idstr

import (
	"errors"

	"github.com/sqids/sqids-go"
)

var s = must(sqids.New(sqids.Options{
	Alphabet: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
}))

func Decode(str string) (int64, error) {
	n := s.Decode(str)

	if len(n) != 1 {
		return 0, errors.New("could not decode")
	}
	return int64(n[0]), nil
}

func Encode(id int64) (string, error) {
	return s.Encode([]uint64{uint64(id)})
}

func must(s *sqids.Sqids, err error) *sqids.Sqids {
	if err != nil {
		panic(err)
	}
	return s
}
