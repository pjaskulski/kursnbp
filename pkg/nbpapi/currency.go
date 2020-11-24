// 'currency' subcommand support - particular currency exchange rates

package nbpapi

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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

// GetCurrency - main function for currrency, selects
// a data download variant depending on previously
// verified input parameters (--table, --code, --date or --last)
func (c *NBPCurrency) GetCurrency(dFlag string, lFlag int, cFlag string, repFormat string) error {
	var err error

	if lFlag != 0 {
		c.result, err = getCurrencyLast(c.tableType, strconv.Itoa(lFlag), cFlag, repFormat)
	} else if dFlag == "today" {
		c.result, err = getCurrencyToday(c.tableType, cFlag, repFormat)
	} else if dFlag == "current" {
		c.result, err = getCurrencyCurrent(c.tableType, cFlag, repFormat)
	} else if len(dFlag) == 10 {
		c.result, err = getCurrencyDay(c.tableType, dFlag, cFlag, repFormat)
	} else if len(dFlag) == 21 {
		c.result, err = getCurrencyRange(c.tableType, dFlag, cFlag, repFormat)
	}
	if err != nil {
		log.Fatal(err)
	}

	if repFormat != "xml" && repFormat != "json" {
		if c.tableType != "C" {
			err = json.Unmarshal(c.result, &c.exchange)
		} else {
			err = json.Unmarshal(c.result, &c.exchangeC)
		}
		if err != nil {
			log.Fatal(err)
		}
	}

	return err
}

// getCurrencyLast - function returns last <last> currency exchange
// rates in json form, or error
func getCurrencyLast(tableType string, last string, currency string, repFormat string) ([]byte, error) {
	address := fmt.Sprintf("%s/rates/%s/%s/last/%s/?format=%s", baseAddressCurrency, tableType, currency, last, repFormat)
	return getJSON(address)
}

// getCurrencyToday - function returns today's currency exchange rate
// in json form, or error
func getCurrencyToday(tableType string, currency string, repFormat string) ([]byte, error) {
	address := fmt.Sprintf("%s/rates/%s/%s/today/?format=%s", baseAddressCurrency, tableType, currency, repFormat)
	return getJSON(address)
}

// getCurrencyCurrent - function returns current exchange rate for
// particular currency (last published price) in json form, or error
func getCurrencyCurrent(tableType string, currency string, repFormat string) ([]byte, error) {
	address := fmt.Sprintf("%s/rates/%s/%s/?format=%s", baseAddressCurrency, tableType, currency, repFormat)
	return getJSON(address)
}

// getCurrencyDay - function returns exchange rate for particular currency
// on the given date (YYYY-MM-DD) in json form, or error
func getCurrencyDay(tableType string, day string, currency string, repFormat string) ([]byte, error) {
	address := fmt.Sprintf("%s/rates/%s/%s/%s/?format=%s", baseAddressCurrency, tableType, currency, day, repFormat)
	return getJSON(address)
}

// getCurrencyRange - function returns exchange rate for particular currency
// within the given date range (RRRR-MM-DD:RRRR-MM-DD) in json form, or error
func getCurrencyRange(tableType string, day string, currency string, repFormat string) ([]byte, error) {
	var startDate string
	var stopDate string

	temp := strings.Split(day, ":")
	startDate = temp[0]
	stopDate = temp[1]

	address := fmt.Sprintf("%s/rates/%s/%s/%s/%s/?format=%s", baseAddressCurrency, tableType, currency, startDate, stopDate, repFormat)
	return getJSON(address)
}

// PrintCurrency - function prints exchange rates as formatted table
// in the console window depending on the tableType field:
// for type A and B tables a column with an average rate is printed,
// for type C two columns: buy price and sell price
func (c *NBPCurrency) PrintCurrency() {
	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.Debug)

	if c.tableType != "C" {
		fmt.Println()
		fmt.Println(l.Get("Table type:")+"\t", c.exchange.Table)
		fmt.Println(l.Get("Currency name:")+"\t", c.exchange.Currency)
		fmt.Println(l.Get("Currency code:")+"\t", c.exchange.Code)
		fmt.Println()

		fmt.Fprintln(w, l.Get("TABLE \t DATE \t AVERAGE (PLN)"))
		fmt.Fprintln(w, l.Get("----- \t ---- \t -------"))
		for _, currencyItem := range c.exchange.Rates {
			currencyValue := fmt.Sprintf("%.4f", currencyItem.Mid)
			fmt.Fprintln(w, currencyItem.No+" \t "+currencyItem.EffectiveDate+" \t "+currencyValue)
		}
	} else {
		fmt.Println()
		fmt.Println(l.Get("Table type:")+"\t", c.exchangeC.Table)
		fmt.Println(l.Get("Currency name:")+"\t", c.exchangeC.Currency)
		fmt.Println(l.Get("Currency code:")+"\t", c.exchangeC.Code)
		fmt.Println()

		fmt.Fprintln(w, l.Get("TABLE \t DATE \t BUY (PLN) \t SELL (PLN) "))
		fmt.Fprintln(w, l.Get("----- \t ---- \t --- \t ---- "))
		for _, currencyItem := range c.exchangeC.Rates {
			currencyValueBid := fmt.Sprintf("%.4f", currencyItem.Bid)
			currencyValueAsk := fmt.Sprintf("%.4f", currencyItem.Ask)
			fmt.Fprintln(w, currencyItem.No+" \t "+currencyItem.EffectiveDate+" \t "+currencyValueBid+" \t "+currencyValueAsk)
		}
	}
	w.Flush()
	fmt.Println()
}

// PrintCurrencyCSV - function prints currency rates in the console,
// in the form of CSV (data separated by a comma), depending on the
// tableType field: for type A and B tables a column with an average
// rate is printed, for type C two columns: buy price and sell price
func (c *NBPCurrency) PrintCurrencyCSV() {
	if c.tableType != "C" {
		fmt.Println(l.Get("TABLE,DATE,AVERAGE (PLN)"))
		for _, currencyItem := range c.exchange.Rates {
			currencyValue := fmt.Sprintf("%.4f", currencyItem.Mid)
			fmt.Println(currencyItem.No + "," + currencyItem.EffectiveDate + "," + currencyValue)
		}
	} else {
		fmt.Println(l.Get("TABLE,DATE,BUY (PLN),SELL (PLN)"))
		for _, currencyItem := range c.exchangeC.Rates {
			currencyValueBid := fmt.Sprintf("%.4f", currencyItem.Bid)
			currencyValueAsk := fmt.Sprintf("%.4f", currencyItem.Ask)
			fmt.Println(currencyItem.No + "," + currencyItem.EffectiveDate + "," + currencyValueBid + "," + currencyValueAsk)
		}
	}
	fmt.Println()
}

// PrintResult - function print just result of request (json or xml)
func (c *NBPCurrency) PrintResult() {
	fmt.Println(string(c.result))
}
