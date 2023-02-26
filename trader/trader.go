package trader

import (
	"errors"
	"fmt"
	"goirc/fin"
	"goirc/model"
	"math"
	"regexp"
	"strconv"
	"strings"
)

func Trade(nick string, msg string) (string, error) {
	re := regexp.MustCompile("^(buy|sell) ([A-Za-z-]+) ([0-9]+)$")
	matches := re.FindStringSubmatch(msg)

	if len(matches) == 0 {
		return "usage: buy|sell symbol shares", nil
	}

	command := strings.ToUpper(string(matches[1]))
	symbol := strings.ToUpper(matches[2])
	amount, err := strconv.Atoi(string(matches[3]))
	if err != nil {
		return "", errors.New(fmt.Sprintf("cannot convert amount: %s", err))
	}

	if amount <= 0 {
		return "", errors.New(fmt.Sprintf("must trade positive amount"))
	}

	if command == "BUY" || command == "SELL" {
		price, err := stockPrice(symbol)
		if err != nil {
			return "", errors.New(fmt.Sprintf("stock price lookup failed: %s", err))
		}

		if command == "BUY" { //////////////////////////////////////////////////////////////// BUY
			cash, err := getCash(nick)
			if err != nil {
				return "", errors.New(fmt.Sprintf("failed lookup cash: %s", err))
			}

			if price*amount > cash {
				return fmt.Sprintf("%d shares of %s costs $%.02f, you have $%.02f", amount, symbol, float64(price*amount)/100.0, float64(cash)/100.0), nil
			}

			err = logTransaction(nick, "BUY", symbol, amount, price)
			if err != nil {
				return "", errors.New(fmt.Sprintf("could not log transaction: %s", err))
			}

			position := Position{
				Symbol: symbol,
				Amount: amount,
				Price:  price,
			}

			report, err := Report(nick)
			if err != nil {
				return "", err
			}

			return fmt.Sprintf("BUY %s | %s", position.String(), report), nil
		}

		if command == "SELL" { //////////////////////////////////////////////////////////////// SELL
			holdings, err := getHoldings(nick)
			if err != nil {
				return "", errors.New(fmt.Sprintf("failed to lookup holdings: %s", err))
			}
			held := (*holdings)[symbol]

			if held < amount {
				return fmt.Sprintf("cannot sell %d %s, holding %d", amount, symbol, held), nil
			}
			err = logTransaction(nick, "SELL", symbol, amount, price)
			if err != nil {
				return "", errors.New(fmt.Sprintf("could not log transaction: %s", err))
			}

			position := Position{
				Symbol: symbol,
				Amount: amount,
				Price:  price,
			}

			report, err := Report(nick)
			if err != nil {
				return "", err
			}

			return fmt.Sprintf("SELL %s | %s", position.String(), report), nil
		}
	}

	return "", errors.New(fmt.Sprintf("unknown command: %s", command))
}

func Report(nick string) (string, error) {
	cash, err := getCash(nick)
	if err != nil {
		return "", err
	}

	holdings, err := getHoldings(nick)
	if err != nil {
		return "", err
	}

	positions, err := holdingsToPositions(*holdings)
	if err != nil {
		return "", err
	}

	arr := []string{
		fmt.Sprintf("CASH $%.02f", float64(cash)/100.0),
	}
	total := cash
	for _, position := range positions {
		arr = append(arr, position.String())
		total += position.Amount * position.Price
	}

	return fmt.Sprintf("%s = NET $%0.2f", strings.Join(arr, " + "), float64(total)/100.0), nil
}

func stockPrice(symbol string) (int, error) {
	resp, err := fin.YahooFinanceFetch(string(symbol))
	if err != nil {
		return 0.0, err
	}

	raw := resp.QuoteSummary.Result[0].FinancialData.CurrentPrice.Raw

	return int(math.Round(raw * 100)), nil
}

type Holdings map[string]int

type Position struct {
	Symbol string
	Amount int
	Price  int
}

func (p *Position) String() string {
	return fmt.Sprintf("%s %d@%.02f $%.02f", p.Symbol, p.Amount, float64(p.Price)/100.0, float64(p.Amount*p.Price)/100.0)
}

type Positions map[string]*Position

func holdingsToPositions(holdings Holdings) (Positions, error) {
	positions := Positions{}

	for symbol, amount := range holdings {
		price, err := stockPrice(symbol)
		if err != nil {
			return nil, err
		}

		if amount > 0 {
			positions[symbol] = &Position{
				Symbol: symbol,
				Amount: amount,
				Price:  price,
			}
		}
	}

	return positions, nil
}

type Transaction struct {
	CreatedAt string
	Nick      string
	Verb      string
	Symbol    string
	Shares    int
	Price     int
}

func getHoldings(nick string) (*Holdings, error) {
	transactions := []Transaction{}
	err := model.DB.Select(&transactions, `select verb, symbol, shares, price from transactions where nick = ? order by created_at asc`, nick)
	if err != nil {
		return nil, err
	}

	holdings := Holdings{}
	for _, tx := range transactions {
		if tx.Verb == "BUY" {
			holdings[tx.Symbol] += tx.Shares
		}
		if tx.Verb == "SELL" {
			holdings[tx.Symbol] -= tx.Shares
		}
	}

	return &holdings, nil
}

func getCash(nick string) (int, error) {
	transactions := []Transaction{}
	err := model.DB.Select(&transactions, `select verb, symbol, shares, price from transactions where nick = ? order by created_at asc`, nick)
	if err != nil {
		return 0, err
	}

	cash := 1_000_000 * 100 // opening balance

	for _, tx := range transactions {
		if tx.Verb == "BUY" {
			cash -= tx.Shares * tx.Price
		}
		if tx.Verb == "SELL" {
			cash += tx.Shares * tx.Price
		}
	}

	return cash, nil
}

func logTransaction(nick string, verb string, symbol string, amount int, price int) error {
	_, err := model.DB.Exec(`insert into transactions(nick, verb, symbol, shares, price) values(?, ?, ?, ?, ?)`, nick, verb, symbol, amount, price)
	return err
}
