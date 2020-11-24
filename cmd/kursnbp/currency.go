// 'currency' subcommand support - particular currency exchange rates

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

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

// getCurrency - main function for currrency, selects
// a data download variant depending on previously
// verified input parameters (--table, --code, --date or --last)
func getCurrency(tFlag string, dFlag string, lFlag int, cFlag string) ([]byte, error) {
	var result []byte
	var err error

	if lFlag != 0 {
		result, err = getCurrencyLast(tFlag, strconv.Itoa(lFlag), cFlag)
	} else if dFlag == "today" {
		result, err = getCurrencyToday(tFlag, cFlag)
	} else if dFlag == "current" {
		result, err = getCurrencyCurrent(tFlag, cFlag)
	} else if len(dFlag) == 10 {
		result, err = getCurrencyDay(tFlag, dFlag, cFlag)
	} else if len(dFlag) == 21 {
		result, err = getCurrencyRange(tFlag, dFlag, cFlag)
	}

	return result, err
}

// getCurrencyLast - function returns last <last> currency exchange
// rates in json form, or error
func getCurrencyLast(tableType string, last string, currency string) ([]byte, error) {
	address := fmt.Sprintf("%s/rates/%s/%s/last/%s/?format=%s", baseAddress, tableType, currency, last, repFormat)
	return getJSON(address)
}

// getCurrencyToday - function returns today's currency exchange rate
// in json form, or error
func getCurrencyToday(tableType string, currency string) ([]byte, error) {
	address := fmt.Sprintf("%s/rates/%s/%s/today/?format=%s", baseAddress, tableType, currency, repFormat)
	return getJSON(address)
}

// getCurrencyCurrent - function returns current exchange rate for
// particular currency (last published price) in json form, or error
func getCurrencyCurrent(tableType string, currency string) ([]byte, error) {
	address := fmt.Sprintf("%s/rates/%s/%s/?format=%s", baseAddress, tableType, currency, repFormat)
	return getJSON(address)
}

// getCurrencyDay - function returns exchange rate for particular currency
// on the given date (YYYY-MM-DD) in json form, or error
func getCurrencyDay(tableType string, day string, currency string) ([]byte, error) {
	address := fmt.Sprintf("%s/rates/%s/%s/%s/?format=%s", baseAddress, tableType, currency, day, repFormat)
	return getJSON(address)
}

// getCurrencyRange - function returns exchange rate for particular currency
// within the given date range (RRRR-MM-DD:RRRR-MM-DD) in json form, or error
func getCurrencyRange(tableType string, day string, currency string) ([]byte, error) {
	var startDate string
	var stopDate string

	temp := strings.Split(day, ":")
	startDate = temp[0]
	stopDate = temp[1]

	address := fmt.Sprintf("%s/rates/%s/%s/%s/%s/?format=%s", baseAddress, tableType, currency, startDate, stopDate, repFormat)
	return getJSON(address)
}

// printCurrency - function prints exchange rates as formatted table
// in the console window
// depending on the tableType parameter: for type A and B tables
// a column with an average rate is printed, for type C two columns:
// buy price and sell price
func printCurrency(result []byte, tableType string) {
	var nbpCurrency exchangeCurrency
	var nbpCurrencyC exchangeCurrencyC

	fmt.Println(appName, "-", appDesc)

	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.Debug)

	if tableType != "C" {
		err := json.Unmarshal(result, &nbpCurrency)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println()
		fmt.Println(l.Get("Table type:")+"\t", nbpCurrency.Table)
		fmt.Println(l.Get("Currency name:")+"\t", nbpCurrency.Currency)
		fmt.Println(l.Get("Currency code:")+"\t", nbpCurrency.Code)
		fmt.Println()

		fmt.Fprintln(w, l.Get("TABLE \t DATE \t AVERAGE (PLN)"))
		fmt.Fprintln(w, l.Get("----- \t ---- \t -------"))
		for _, currencyItem := range nbpCurrency.Rates {
			currencyValue := fmt.Sprintf("%.4f", currencyItem.Mid)
			fmt.Fprintln(w, currencyItem.No+" \t "+currencyItem.EffectiveDate+" \t "+currencyValue)
		}
	} else {
		err := json.Unmarshal(result, &nbpCurrencyC)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println()
		fmt.Println(l.Get("Table type:")+"\t", nbpCurrencyC.Table)
		fmt.Println(l.Get("Currency name:")+"\t", nbpCurrencyC.Currency)
		fmt.Println(l.Get("Currency code:")+"\t", nbpCurrencyC.Code)
		fmt.Println()

		fmt.Fprintln(w, l.Get("TABLE \t DATE \t BUY (PLN) \t SELL (PLN) "))
		fmt.Fprintln(w, l.Get("----- \t ---- \t --- \t ---- "))
		for _, currencyItem := range nbpCurrencyC.Rates {
			currencyValueBid := fmt.Sprintf("%.4f", currencyItem.Bid)
			currencyValueAsk := fmt.Sprintf("%.4f", currencyItem.Ask)
			fmt.Fprintln(w, currencyItem.No+" \t "+currencyItem.EffectiveDate+" \t "+currencyValueBid+" \t "+currencyValueAsk)
		}
	}
	w.Flush()

	fmt.Println()
}

// printCurrencyCSV - function prints currency rates in the console,
// in the form of CSV (data separated by a comma), depending on the
// tableType parameter: for type A and B tables a column with an average
// rate is printed, for type C two columns: buy price and sell price
func printCurrencyCSV(result []byte, tableType string) {
	var nbpCurrency exchangeCurrency
	var nbpCurrencyC exchangeCurrencyC

	if tableType != "C" {
		err := json.Unmarshal(result, &nbpCurrency)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(l.Get("TABLE,DATE,AVERAGE (PLN)"))
		for _, currencyItem := range nbpCurrency.Rates {
			currencyValue := fmt.Sprintf("%.4f", currencyItem.Mid)
			fmt.Println(currencyItem.No + "," + currencyItem.EffectiveDate + "," + currencyValue)
		}
	} else {
		err := json.Unmarshal(result, &nbpCurrencyC)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(l.Get("TABLE,DATE,BUY (PLN),SELL (PLN)"))
		for _, currencyItem := range nbpCurrencyC.Rates {
			currencyValueBid := fmt.Sprintf("%.4f", currencyItem.Bid)
			currencyValueAsk := fmt.Sprintf("%.4f", currencyItem.Ask)
			fmt.Println(currencyItem.No + "," + currencyItem.EffectiveDate + "," + currencyValueBid + "," + currencyValueAsk)
		}
	}
	fmt.Println()
}
