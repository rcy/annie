package gold

import (
	"os"
	"testing"
)

func TestGetGoldPrice(t *testing.T) {
	var token = os.Getenv("GOLD_API_TOKEN")
	price, err := getGoldPrice(token)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Price: %f\n", price)
}
