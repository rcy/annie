package idstr

import (
	"errors"

	"github.com/sqids/sqids-go"
)

func Decode(str string) (int64, error) {
	s, _ := sqids.New()

	n := s.Decode(str)

	if len(n) != 1 {
		return 0, errors.New("could not decode")
	}
	return int64(n[0]), nil
}

func Encode(id int64) (string, error) {
	s, _ := sqids.New()

	return s.Encode([]uint64{uint64(id)})
}
