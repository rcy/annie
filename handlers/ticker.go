package handlers

import (
	"goirc/bot"
	"goirc/fin"
	"goirc/util"
	"strings"
)

func Ticker(params bot.HandlerParams) bool {
	symbol := params.Matches[1]

	data, err := fin.YahooFinanceFetch(symbol)
	if err != nil {
		params.Privmsgf(params.Target, "error: %s", err)
		return true
	}

	result := data.QuoteSummary.Result[0]
	params.Privmsgf(params.Target, "%s %s %f", strings.ToUpper(symbol), util.BareDomain(result.SummaryProfile.Website), result.FinancialData.CurrentPrice.Raw)

	return true
}
