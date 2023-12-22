package idstr

import (
	"fmt"
	"testing"
)

var tt = []struct {
	num int64
	str string
}{
	{
		num: 1,
		str: "Uk",
	},
	{
		num: 2,
		str: "gb",
	},
	{
		num: 69,
		str: "ARb",
	},
	{
		num: 1000,
		str: "pnd",
	},
	{
		num: 23443,
		str: "AV7d",
	},
	{
		num: 999999,
		str: "UQ1Nd",
	},
}

func TestEncode(t *testing.T) {
	for _, tc := range tt {
		t.Run(fmt.Sprint(tc.num), func(t *testing.T) {
			got, _ := Encode(tc.num)
			if got != tc.str {
				t.Errorf("expected %s got %s", tc.str, got)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	for _, tc := range tt {
		t.Run(tc.str, func(t *testing.T) {
			got, _ := Decode(tc.str)
			if got != tc.num {
				t.Errorf("expected %d got %d", tc.num, got)
			}
		})
	}
}
