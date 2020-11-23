// 'table' subcommand support - complete tables of currency exchange rates

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

type rateTable struct {
	Currency string  `json:"currency"`
	Code     string  `json:"code"`
	Mid      float64 `json:"mid"`
}

type exchangeTable struct {
	Table         string      `json:"table"`
	No            string      `json:"no"`
	EffectiveDate string      `json:"effectiveDate"`
	Rates         []rateTable `json:"rates"`
}

type rateTableC struct {
	Currency string  `json:"currency"`
	Code     string  `json:"code"`
	Bid      float64 `json:"bid"`
	Ask      float64 `json:"ask"`
}

type exchangeTableC struct {
	Table         string       `json:"table"`
	No            string       `json:"no"`
	TradingDate   string       `json:"tradingDate"`
	EffectiveDate string       `json:"effectiveDate"`
	Rates         []rateTableC `json:"rates"`
}

// getTable - main download function for table, selects
// a data download variant depending on previously
// verified input parameters (--table, --date or --last)
func getTable(tFlag string, dFlag string, lFlag int) ([]byte, error) {
	var result []byte
	var err error

	if lFlag != 0 {
		result, err = getTableLast(tFlag, strconv.Itoa(lFlag))
	} else if dFlag == "today" {
		result, err = getTableToday(tFlag)
	} else if dFlag == "current" {
		result, err = getTableCurrent(tFlag)
	} else if len(dFlag) == 10 {
		result, err = getTableDay(tFlag, dFlag)
	} else if len(dFlag) == 21 {
		result, err = getTableRange(tFlag, dFlag)
	}

	return result, err
}

// getTableToday - function returns exchange rate table published today
// in JSON form, or error
func getTableToday(tableType string) ([]byte, error) {
	address := fmt.Sprintf(baseAddress+"/tables/%s/today/?format="+repFormat, tableType)
	return getJSON(address)
}

// getTableCurrent - function returns current table of exchange rates
// (last published table) in JSON form, or error
func getTableCurrent(tableType string) ([]byte, error) {
	address := fmt.Sprintf(baseAddress+"/tables/%s/?format="+repFormat, tableType)
	return getJSON(address)
}

// getTableDay - functions returns table of exchange rates
// on the given date (YYYY-MM-DD) in JSON form, or error
func getTableDay(tableType string, day string) ([]byte, error) {
	address := fmt.Sprintf(baseAddress+"/tables/%s/%s/?format="+repFormat, tableType, day)
	return getJSON(address)
}

// getTableRange - function returns table of exchange rates  within
// the given date range (RRRR-MM-DD:RRRR-MM-DD) in JSON form, or error
func getTableRange(tableType string, day string) ([]byte, error) {
	var startDate string
	var stopDate string

	temp := strings.Split(day, ":")
	if len(temp) != 2 {
		log.Fatal(errors.New(l.Get("Invalid date range format")))
	}

	startDate = temp[0]
	stopDate = temp[1]

	address := fmt.Sprintf(baseAddress+"/tables/%s/%s/%s/?format="+repFormat, tableType, startDate, stopDate)
	return getJSON(address)
}

// getTableLast - function returns last <last> tables of exchange rates
// in JSON form, or error
func getTableLast(tableType string, last string) ([]byte, error) {
	address := fmt.Sprintf(baseAddress+"/tables/%s/last/%s/?format="+repFormat, tableType, last)
	return getJSON(address)
}

// printTable - function prints tables of exchange rates as
// formatted table in the console window,
// depending on the tableType parameter: for type A and B tables
// a column with an average rate is printed, for type C two columns:
// buy price and sell price
func printTable(result []byte, tableType string) {
	var nbpTables []exchangeTable
	var nbpTablesC []exchangeTableC

	fmt.Println(appName, "-", appDesc)

	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.Debug)

	if tableType != "C" {
		err := json.Unmarshal(result, &nbpTables)
		if err != nil {
			log.Fatal(err)
		}

		for _, item := range nbpTables {
			fmt.Println()
			fmt.Println(l.Get("Table type:")+"\t\t", item.Table)
			fmt.Println(l.Get("Table number:")+"\t\t", item.No)
			fmt.Println(l.Get("Publication date:")+"\t", item.EffectiveDate)
			fmt.Println()

			fmt.Fprintln(w, l.Get("CODE \t NAME \t AVERAGE (PLN)"))
			fmt.Fprintln(w, l.Get("---- \t ---- \t -------------"))
			for _, currencyItem := range item.Rates {
				currencyValue := fmt.Sprintf("%.4f", currencyItem.Mid)
				fmt.Fprintln(w, currencyItem.Code+" \t "+currencyItem.Currency+" \t "+currencyValue)
			}

			w.Flush()
		}
	} else {
		err := json.Unmarshal(result, &nbpTablesC)
		if err != nil {
			log.Fatal(err)
		}

		for _, item := range nbpTablesC {
			fmt.Println()
			fmt.Println(l.Get("Table type:")+"\t\t", item.Table)
			fmt.Println(l.Get("Table number:")+"\t\t", item.No)
			fmt.Println(l.Get("Trading date:")+"\t\t", item.TradingDate)
			fmt.Println(l.Get("Publication date:")+"\t", item.EffectiveDate)
			fmt.Println()

			fmt.Fprintln(w, l.Get("CODE \t NAME \t BUY (PLN) \t SELL (PLN) "))
			fmt.Fprintln(w, l.Get("---- \t ---- \t --------- \t ---------- "))
			for _, currencyItem := range item.Rates {
				currencyValueBid := fmt.Sprintf("%.4f", currencyItem.Bid)
				currencyValueAsk := fmt.Sprintf("%.4f", currencyItem.Ask)
				fmt.Fprintln(w, currencyItem.Code+" \t "+currencyItem.Currency+" \t "+currencyValueBid+" \t "+currencyValueAsk)
			}

			w.Flush()
		}
	}

	fmt.Println()
}

// printTableCSV - function prints tables of exchange rates in the console,
// in the form of CSV (data separated by a comma), depending on the
// tableType parameter: for type A and B tables a column with an average
// rate is printed, for type C two columns: buy price and sell price
func printTableCSV(result []byte, tableType string) {
	var nbpTables []exchangeTable
	var nbpTablesC []exchangeTableC
	var tableNo string

	if tableType != "C" {
		err := json.Unmarshal(result, &nbpTables)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(l.Get("TABLE,CODE,NAME,AVERAGE (PLN)"))

		for _, item := range nbpTables {
			tableNo = item.No
			for _, currencyItem := range item.Rates {
				currencyValue := fmt.Sprintf("%.4f", currencyItem.Mid)
				fmt.Println(tableNo + "," + currencyItem.Code + "," + currencyItem.Currency + "," + currencyValue)
			}
		}
	} else {
		err := json.Unmarshal(result, &nbpTablesC)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(l.Get("TABLE,CODE,NAME,BUY (PLN),SELL (PLN)"))

		for _, item := range nbpTablesC {
			tableNo = item.No
			for _, currencyItem := range item.Rates {
				currencyValueBid := fmt.Sprintf("%.4f", currencyItem.Bid)
				currencyValueAsk := fmt.Sprintf("%.4f", currencyItem.Ask)
				fmt.Println(tableNo + "," + currencyItem.Code + "," + currencyItem.Currency + "," + currencyValueBid + "," + currencyValueAsk)
			}
			fmt.Println()
		}
	}

	fmt.Println()
}
