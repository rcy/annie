package fin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type YahooResponse struct {
	QuoteSummary struct {
		Result []struct {
			SummaryProfile struct {
				Address1            string        `json:"address1"`
				Address2            string        `json:"address2"`
				City                string        `json:"city"`
				State               string        `json:"state"`
				Zip                 string        `json:"zip"`
				Country             string        `json:"country"`
				Phone               string        `json:"phone"`
				Website             string        `json:"website"`
				Industry            string        `json:"industry"`
				Sector              string        `json:"sector"`
				LongBusinessSummary string        `json:"longBusinessSummary"`
				FullTimeEmployees   int           `json:"fullTimeEmployees"`
				CompanyOfficers     []interface{} `json:"companyOfficers"`
				MaxAge              int           `json:"maxAge"`
			} `json:"summaryProfile"`

			FinancialData struct {
				MaxAge       int `json:"maxAge"`
				CurrentPrice struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"currentPrice"`
				TargetHighPrice struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"targetHighPrice"`
				TargetLowPrice struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"targetLowPrice"`
				TargetMeanPrice struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"targetMeanPrice"`
				TargetMedianPrice struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"targetMedianPrice"`
				RecommendationMean struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"recommendationMean"`
				RecommendationKey       string `json:"recommendationKey"`
				NumberOfAnalystOpinions struct {
					Raw     int    `json:"raw"`
					Fmt     string `json:"fmt"`
					LongFmt string `json:"longFmt"`
				} `json:"numberOfAnalystOpinions"`
				TotalCash struct {
					Raw     int64  `json:"raw"`
					Fmt     string `json:"fmt"`
					LongFmt string `json:"longFmt"`
				} `json:"totalCash"`
				TotalCashPerShare struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"totalCashPerShare"`
				Ebitda struct {
					Raw     int64  `json:"raw"`
					Fmt     string `json:"fmt"`
					LongFmt string `json:"longFmt"`
				} `json:"ebitda"`
				TotalDebt struct {
					Raw     int64  `json:"raw"`
					Fmt     string `json:"fmt"`
					LongFmt string `json:"longFmt"`
				} `json:"totalDebt"`
				QuickRatio struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"quickRatio"`
				CurrentRatio struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"currentRatio"`
				TotalRevenue struct {
					Raw     int64  `json:"raw"`
					Fmt     string `json:"fmt"`
					LongFmt string `json:"longFmt"`
				} `json:"totalRevenue"`
				DebtToEquity struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"debtToEquity"`
				RevenuePerShare struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"revenuePerShare"`
				ReturnOnAssets struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"returnOnAssets"`
				ReturnOnEquity struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"returnOnEquity"`
				GrossProfits struct {
					Raw     int64  `json:"raw"`
					Fmt     string `json:"fmt"`
					LongFmt string `json:"longFmt"`
				} `json:"grossProfits"`
				FreeCashflow struct {
					Raw     int64  `json:"raw"`
					Fmt     string `json:"fmt"`
					LongFmt string `json:"longFmt"`
				} `json:"freeCashflow"`
				OperatingCashflow struct {
					Raw     int64  `json:"raw"`
					Fmt     string `json:"fmt"`
					LongFmt string `json:"longFmt"`
				} `json:"operatingCashflow"`
				EarningsGrowth struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"earningsGrowth"`
				RevenueGrowth struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"revenueGrowth"`
				GrossMargins struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"grossMargins"`
				EbitdaMargins struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"ebitdaMargins"`
				OperatingMargins struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"operatingMargins"`
				ProfitMargins struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"profitMargins"`
				FinancialCurrency string `json:"financialCurrency"`
			} `json:"financialData"`
		} `json:"result"`
		Error struct {
			Code        string `json:"code"`
			Description string `json:"description"`
		} `json:"error"`
	} `json:"quoteSummary"`
}

func YahooFinanceFetch(symbol string) (*YahooResponse, error) {
	//curl 'https://query1.finance.yahoo.com/v11/finance/quoteSummary/tsla?modules=financialData'|jq .
	resp, err := http.Get(fmt.Sprintf("https://query1.finance.yahoo.com/v11/finance/quoteSummary/%s?modules=financialData,summaryProfile", symbol))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := YahooResponse{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	if data.QuoteSummary.Error.Code != "" {
		return nil, errors.New(data.QuoteSummary.Error.Code)
	}

	return &data, nil
}
