package handlers

import (
	"goirc/fin"
	"goirc/util"
	"regexp"
	"strings"
)

func ticker(params Params) bool {
	re := regexp.MustCompile("^[$]([A-Za-z-]+)")
	matches := re.FindSubmatch([]byte(params.Msg))

	if len(matches) == 0 {
		return false
	}

	symbol := string(matches[1])

	data, err := fin.YahooFinanceFetch(symbol)
	if err != nil {
		params.Privmsgf(params.Target, "error: %s", err)
		return true
	}

	result := data.QuoteSummary.Result[0]
	params.Privmsgf(params.Target, "%s %s %f", strings.ToUpper(symbol), util.BareDomain(result.SummaryProfile.Website), result.FinancialData.CurrentPrice.Raw)

	return true
}
