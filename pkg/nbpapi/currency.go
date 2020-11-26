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
	exchange  exchangeCurrency
	exchangeC exchangeCurrencyC
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

// GetCurrencyRaw - function downloads data in json or xml form
func (c *NBPCurrency) GetCurrencyRaw(dFlag string, lFlag int, cFlag string, repFormat string) error {
	var err error

	address := getCurrencyAddress(c.tableType, dFlag, lFlag, cFlag)
	c.result, err = getData(address, repFormat)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

// GetCurrencyByDate - function downloads and writes data to exchange (exchangeC) slice,
// raw data (json) still available in result field
func (c *NBPCurrency) GetCurrencyByDate(dFlag string, cFlag string) error {
	var err error

	address := getCurrencyAddress(c.tableType, dFlag, 0, cFlag)
	c.result, err = getData(address, "json")
	if err != nil {
		log.Fatal(err)
	}

	if c.tableType != "C" {
		err = json.Unmarshal(c.result, &c.exchange)
	} else {
		err = json.Unmarshal(c.result, &c.exchangeC)
	}
	if err != nil {
		log.Fatal(err)
	}

	return err
}

// GetCurrencyLast - function downloads and writes data to exchange (exchangeC) slice,
// raw data (json) still available in result field
func (c *NBPCurrency) GetCurrencyLast(lFlag int, cFlag string) error {
	var err error

	address := getCurrencyAddress(c.tableType, "", lFlag, cFlag)
	c.result, err = getData(address, "json")
	if err != nil {
		log.Fatal(err)
	}

	if c.tableType != "C" {
		err = json.Unmarshal(c.result, &c.exchange)
	} else {
		err = json.Unmarshal(c.result, &c.exchangeC)
	}
	if err != nil {
		log.Fatal(err)
	}

	return err
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
		output += fmt.Sprintln(l.Get("Table type:")+"\t", c.exchange.Table)
		output += fmt.Sprintln(l.Get("Currency name:")+"\t", c.exchange.Currency)
		output += fmt.Sprintln(l.Get("Currency code:")+"\t", c.exchange.Code)
		output += fmt.Sprintln()

		fmt.Fprintln(w, l.Get("TABLE \t DATE \t AVERAGE (PLN)"))
		fmt.Fprintln(w, l.Get("----- \t ---- \t -------------"))
		for _, currencyItem := range c.exchange.Rates {
			currencyValue := fmt.Sprintf("%.4f", currencyItem.Mid)
			fmt.Fprintln(w, currencyItem.No+" \t "+currencyItem.EffectiveDate+" \t "+currencyValue)
		}
	} else {
		output += fmt.Sprintln()
		output += fmt.Sprintln(l.Get("Table type:")+"\t", c.exchangeC.Table)
		output += fmt.Sprintln(l.Get("Currency name:")+"\t", c.exchangeC.Currency)
		output += fmt.Sprintln(l.Get("Currency code:")+"\t", c.exchangeC.Code)
		output += fmt.Sprintln()

		fmt.Fprintln(w, l.Get("TABLE \t DATE \t BUY (PLN) \t SELL (PLN) "))
		fmt.Fprintln(w, l.Get("----- \t ---- \t --------- \t ---------- "))
		for _, currencyItem := range c.exchangeC.Rates {
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
		for _, currencyItem := range c.exchange.Rates {
			currencyValue := fmt.Sprintf("%.4f", currencyItem.Mid)
			output += fmt.Sprintln(currencyItem.No + "," + currencyItem.EffectiveDate + "," + currencyValue)
		}
	} else {
		output += fmt.Sprintln(l.Get("TABLE,DATE,BUY (PLN),SELL (PLN)"))
		for _, currencyItem := range c.exchangeC.Rates {
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
