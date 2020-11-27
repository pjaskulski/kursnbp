// 'currency' subcommand support - particular currency exchange rates

package nbpapi

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"text/tabwriter"
)

// base addresses of the NBP API service
const (
	baseAddressCurrency string = "http://api.nbp.pl/api/exchangerates"
)

// NBPCurrency type
type NBPCurrency struct {
	tableType string
	result    []byte
	Exchange  exchangeCurrency
	ExchangeC exchangeCurrencyC
}

type rateCurrency struct {
	No            string  `json:"no"`
	EffectiveDate string  `json:"effectiveDate"`
	Mid           float64 `json:"mid"`
}

type exchangeCurrency struct {
	Table    string         `json:"table"`
	Currency string         `json:"currency"`
	Code     string         `json:"code"`
	Rates    []rateCurrency `json:"rates"`
}

type rateCurrencyC struct {
	No            string  `json:"no"`
	EffectiveDate string  `json:"effectiveDate"`
	Bid           float64 `json:"bid"`
	Ask           float64 `json:"ask"`
}

type exchangeCurrencyC struct {
	Table    string          `json:"table"`
	Currency string          `json:"currency"`
	Code     string          `json:"code"`
	Rates    []rateCurrencyC `json:"rates"`
}

// Rate type
type Rate struct {
	No            string
	EffectiveDate string
	Mid           float64
	Bid           float64
	Ask           float64
}

// NewCurrency - function creates new currency type
func NewCurrency(tFlag string) *NBPCurrency {
	return &NBPCurrency{
		tableType: tFlag,
	}
}

// getCurrencyAddress - function builds download address depending on previously
// verified input parameters (--table, --date or --last, --code)
func getCurrencyAddress(tableType string, dFlag string, lFlag int, cFlag string) string {
	var address string

	if lFlag != 0 {
		address = queryCurrencyLast(tableType, strconv.Itoa(lFlag), cFlag)
	} else if dFlag == "today" {
		address = queryCurrencyToday(tableType, cFlag)
	} else if dFlag == "current" {
		address = queryCurrencyCurrent(tableType, cFlag)
	} else if len(dFlag) == 10 {
		address = queryCurrencyDay(tableType, dFlag, cFlag)
	} else if len(dFlag) == 21 {
		address = queryCurrencyRange(tableType, dFlag, cFlag)
	}

	return address
}

// CurrencyRaw - function downloads data in json or xml form
func (c *NBPCurrency) CurrencyRaw(dFlag string, lFlag int, cFlag string, repFormat string) error {
	var err error

	address := getCurrencyAddress(c.tableType, dFlag, lFlag, cFlag)
	c.result, err = getData(address, repFormat)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

// CurrencyByDate - function downloads and writes data to exchange (exchangeC) slice,
// raw data (json) still available in result field
func (c *NBPCurrency) CurrencyByDate(dFlag string, cFlag string) error {
	var err error

	address := getCurrencyAddress(c.tableType, dFlag, 0, cFlag)
	c.result, err = getData(address, "json")
	if err != nil {
		log.Fatal(err)
	}

	if c.tableType != "C" {
		err = json.Unmarshal(c.result, &c.Exchange)
	} else {
		err = json.Unmarshal(c.result, &c.ExchangeC)
	}
	if err != nil {
		log.Fatal(err)
	}

	return err
}

// CurrencyLast - function downloads and writes data to exchange (exchangeC) slice,
// raw data (json) still available in result field
func (c *NBPCurrency) CurrencyLast(cFlag string, lFlag int) error {
	var err error

	address := getCurrencyAddress(c.tableType, "", lFlag, cFlag)
	c.result, err = getData(address, "json")
	if err != nil {
		log.Fatal(err)
	}

	if c.tableType != "C" {
		err = json.Unmarshal(c.result, &c.Exchange)
	} else {
		err = json.Unmarshal(c.result, &c.ExchangeC)
	}
	if err != nil {
		log.Fatal(err)
	}

	return err
}

// CurrencyToday - function downloads and writes data to exchange (exchangeC) slice,
// raw data (json) still available in result field
func (c *NBPCurrency) CurrencyToday(cFlag string) error {
	var err error

	address := getCurrencyAddress(c.tableType, "today", 0, cFlag)
	c.result, err = getData(address, "json")
	if err != nil {
		log.Fatal(err)
	}

	if c.tableType != "C" {
		err = json.Unmarshal(c.result, &c.Exchange)
	} else {
		err = json.Unmarshal(c.result, &c.ExchangeC)
	}
	if err != nil {
		log.Fatal(err)
	}

	return err
}

// GetRateCurrent - function downloads current currency exchange rate
// and return Rate struct (or error)
func (c *NBPCurrency) GetRateCurrent(cFlag string) (Rate, error) {
	var err error

	address := getCurrencyAddress(c.tableType, "current", 0, cFlag)
	c.result, err = getData(address, "json")
	if err != nil {
		log.Fatal(err)
	}

	if c.tableType != "C" {
		err = json.Unmarshal(c.result, &c.Exchange)
	} else {
		err = json.Unmarshal(c.result, &c.ExchangeC)
	}
	if err != nil {
		log.Fatal(err)
	}

	var rate Rate
	if c.tableType != "C" {
		rate.No = c.Exchange.Rates[0].No
		rate.EffectiveDate = c.Exchange.Rates[0].EffectiveDate
		rate.Mid = c.Exchange.Rates[0].Mid
		rate.Ask = 0
		rate.Bid = 0
	} else {
		rate.No = c.ExchangeC.Rates[0].No
		rate.EffectiveDate = c.ExchangeC.Rates[0].EffectiveDate
		rate.Mid = 0
		rate.Ask = c.ExchangeC.Rates[0].Ask
		rate.Bid = c.ExchangeC.Rates[0].Bid
	}

	return rate, err
}

// GetRateToday - function downloads today's currency exchange rate
// and returns Rate struct (or error)
func (c *NBPCurrency) GetRateToday(cFlag string) (Rate, error) {
	var err error

	address := getCurrencyAddress(c.tableType, "today", 0, cFlag)
	c.result, err = getData(address, "json")
	if err != nil {
		log.Fatal(err)
	}

	if c.tableType != "C" {
		err = json.Unmarshal(c.result, &c.Exchange)
	} else {
		err = json.Unmarshal(c.result, &c.ExchangeC)
	}
	if err != nil {
		log.Fatal(err)
	}

	var rate Rate
	if c.tableType != "C" {
		rate.No = c.Exchange.Rates[0].No
		rate.EffectiveDate = c.Exchange.Rates[0].EffectiveDate
		rate.Mid = c.Exchange.Rates[0].Mid
		rate.Ask = 0
		rate.Bid = 0
	} else {
		rate.No = c.ExchangeC.Rates[0].No
		rate.EffectiveDate = c.ExchangeC.Rates[0].EffectiveDate
		rate.Mid = 0
		rate.Ask = c.ExchangeC.Rates[0].Ask
		rate.Bid = c.ExchangeC.Rates[0].Bid
	}

	return rate, err
}

// GetRateByDate - function downloads today's currency exchange rate
// and returns Rate struct (or error)
func (c *NBPCurrency) GetRateByDate(code string, date string) ([]Rate, error) {
	var err error

	err = CheckArg("currency", c.tableType, date, 0, "table", code)
	if err != nil {
		return nil, err
	}

	address := getCurrencyAddress(c.tableType, date, 0, code)
	c.result, err = getData(address, "json")
	if err != nil {
		return nil, err
	}

	if c.tableType != "C" {
		err = json.Unmarshal(c.result, &c.Exchange)
	} else {
		err = json.Unmarshal(c.result, &c.ExchangeC)
	}
	if err != nil {
		return nil, err
	}

	var rates []Rate
	var rate Rate
	if c.tableType != "C" {
		for _, item := range c.Exchange.Rates {
			rate.No = item.No
			rate.EffectiveDate = item.EffectiveDate
			rate.Mid = item.Mid
			rate.Ask = 0
			rate.Bid = 0
			rates = append(rates, rate)
		}
	} else {
		for _, item := range c.ExchangeC.Rates {
			rate.No = item.No
			rate.EffectiveDate = item.EffectiveDate
			rate.Mid = 0
			rate.Ask = item.Ask
			rate.Bid = item.Bid
			rates = append(rates, rate)
		}
	}

	return rates, err
}

// queryCurrencyLast - returns query: last <number> currency exchange
// rates in json/xml form, or error
func queryCurrencyLast(tableType string, last string, currency string) string {
	return fmt.Sprintf("%s/rates/%s/%s/last/%s/", baseAddressCurrency, tableType, currency, last)

}

// queryCurrencyToday - returns query: today's currency exchange rate
func queryCurrencyToday(tableType string, currency string) string {
	return fmt.Sprintf("%s/rates/%s/%s/today/", baseAddressCurrency, tableType, currency)
}

// queryCurrencyCurrent - returns query: current exchange rate for
// particular currency (last published price)
func queryCurrencyCurrent(tableType string, currency string) string {
	return fmt.Sprintf("%s/rates/%s/%s/", baseAddressCurrency, tableType, currency)
}

// queryCurrencyDay - returns query: exchange rate for particular currency
// on the given date (YYYY-MM-DD)
func queryCurrencyDay(tableType string, day string, currency string) string {
	return fmt.Sprintf("%s/rates/%s/%s/%s/", baseAddressCurrency, tableType, currency, day)
}

// queryCurrencyRange - returns query: exchange rate for particular currency
// within the given date range (RRRR-MM-DD:RRRR-MM-DD)
func queryCurrencyRange(tableType string, day string, currency string) string {
	var startDate string
	var stopDate string

	temp := strings.Split(day, ":")
	startDate = temp[0]
	stopDate = temp[1]

	address := fmt.Sprintf("%s/rates/%s/%s/%s/%s/", baseAddressCurrency, tableType, currency, startDate, stopDate)
	return address
}

// GetPrettyOutput - function returns exchange rates as formatted table
// depending on the tableType field:
// for type A and B tables a column with an average rate is printed,
// for type C two columns: buy price and sell price
func (c *NBPCurrency) GetPrettyOutput() string {
	const padding = 3
	var builder strings.Builder
	var output string
	w := tabwriter.NewWriter(&builder, 0, 0, padding, ' ', tabwriter.Debug)

	if c.tableType != "C" {
		output += fmt.Sprintln()
		output += fmt.Sprintln(l.Get("Table type:")+"\t", c.Exchange.Table)
		output += fmt.Sprintln(l.Get("Currency name:")+"\t", c.Exchange.Currency)
		output += fmt.Sprintln(l.Get("Currency code:")+"\t", c.Exchange.Code)
		output += fmt.Sprintln()

		fmt.Fprintln(w, l.Get("TABLE \t DATE \t AVERAGE (PLN)"))
		fmt.Fprintln(w, l.Get("----- \t ---- \t -------------"))
		for _, currencyItem := range c.Exchange.Rates {
			currencyValue := fmt.Sprintf("%.4f", currencyItem.Mid)
			fmt.Fprintln(w, currencyItem.No+" \t "+currencyItem.EffectiveDate+" \t "+currencyValue)
		}
	} else {
		output += fmt.Sprintln()
		output += fmt.Sprintln(l.Get("Table type:")+"\t", c.ExchangeC.Table)
		output += fmt.Sprintln(l.Get("Currency name:")+"\t", c.ExchangeC.Currency)
		output += fmt.Sprintln(l.Get("Currency code:")+"\t", c.ExchangeC.Code)
		output += fmt.Sprintln()

		fmt.Fprintln(w, l.Get("TABLE \t DATE \t BUY (PLN) \t SELL (PLN) "))
		fmt.Fprintln(w, l.Get("----- \t ---- \t --------- \t ---------- "))
		for _, currencyItem := range c.ExchangeC.Rates {
			currencyValueBid := fmt.Sprintf("%.4f", currencyItem.Bid)
			currencyValueAsk := fmt.Sprintf("%.4f", currencyItem.Ask)
			fmt.Fprintln(w, currencyItem.No+" \t "+currencyItem.EffectiveDate+" \t "+currencyValueBid+" \t "+currencyValueAsk)
		}
	}
	w.Flush()

	return output + builder.String()
}

// GetCSVOutput - function returns currency rates,
// in the form of CSV (data separated by a comma), depending on the
// tableType field: for type A and B tables a column with an average
// rate is printed, for type C two columns: buy price and sell price
func (c *NBPCurrency) GetCSVOutput() string {
	var output string = ""

	if c.tableType != "C" {
		output += fmt.Sprintln(l.Get("TABLE,DATE,AVERAGE (PLN)"))
		for _, currencyItem := range c.Exchange.Rates {
			currencyValue := fmt.Sprintf("%.4f", currencyItem.Mid)
			output += fmt.Sprintln(currencyItem.No + "," + currencyItem.EffectiveDate + "," + currencyValue)
		}
	} else {
		output += fmt.Sprintln(l.Get("TABLE,DATE,BUY (PLN),SELL (PLN)"))
		for _, currencyItem := range c.ExchangeC.Rates {
			currencyValueBid := fmt.Sprintf("%.4f", currencyItem.Bid)
			currencyValueAsk := fmt.Sprintf("%.4f", currencyItem.Ask)
			output += fmt.Sprintln(currencyItem.No + "," + currencyItem.EffectiveDate + "," + currencyValueBid + "," + currencyValueAsk)
		}
	}

	return output
}

// GetRawOutput - function print just result of request (json or xml)
func (c *NBPCurrency) GetRawOutput() string {
	return string(c.result)
}
