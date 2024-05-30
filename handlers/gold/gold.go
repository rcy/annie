package gold

import (
	"encoding/json"
	"fmt"
	"goirc/bot"
	"net/http"
	"os"
)

type response struct {
	Timestamp      int     `json:"timestamp"`
	Metal          string  `json:"metal"`
	Currency       string  `json:"currency"`
	Exchange       string  `json:"exchange"`
	Symbol         string  `json:"symbol"`
	PrevClosePrice float64 `json:"prev_close_price"`
	OpenPrice      float64 `json:"open_price"`
	LowPrice       float64 `json:"low_price"`
	HighPrice      float64 `json:"high_price"`
	OpenTime       int     `json:"open_time"`
	Price          float64 `json:"price"`
	Ch             float64 `json:"ch"`
	Chp            float64 `json:"chp"`
	Ask            float64 `json:"ask"`
	Bid            float64 `json:"bid"`
	PriceGram24K   float64 `json:"price_gram_24k"`
	PriceGram22K   float64 `json:"price_gram_22k"`
	PriceGram21K   float64 `json:"price_gram_21k"`
	PriceGram20K   float64 `json:"price_gram_20k"`
	PriceGram18K   float64 `json:"price_gram_18k"`
	PriceGram16K   float64 `json:"price_gram_16k"`
	PriceGram14K   float64 `json:"price_gram_14k"`
	PriceGram10K   float64 `json:"price_gram_10k"`
}

func getGoldPrice(token string) (float64, error) {
	url := "https://www.goldapi.io/api/XAU/CAD"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0.0, err
	}
	req.Header.Set("x-access-token", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0.0, err
	}
	defer resp.Body.Close()

	var obj response
	json.NewDecoder(resp.Body).Decode(&obj)

	fmt.Printf("Price: %f\n", obj.Price)

	return obj.Price, nil
}

func Handle(params bot.HandlerParams) error {
	price, err := getGoldPrice(os.Getenv("GOLD_API_TOKEN"))
	if err != nil {
		return err
	}

	params.Privmsgf(params.Target, "gold is $%.2f / oz t.", price)

	return nil
}
