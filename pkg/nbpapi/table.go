// 'table' subcommand support - complete tables of currency exchange rates

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
	baseAddressTable string = "http://api.nbp.pl/api/exchangerates"
)

// NBPTable type
type NBPTable struct {
	tableType string
	result    []byte
	exchange  []exchangeTable
	exchangeC []exchangeTableC
}

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

// NewTable - function creates new table type
func NewTable(tFlag string) *NBPTable {
	return &NBPTable{
		tableType: tFlag,
	}
}

// GetTable - main download function for table, selects
// a data download variant depending on previously
// verified input parameters (--table, --date or --last)
func (t *NBPTable) GetTable(dFlag string, lFlag int, repFormat string) error {
	var err error

	if lFlag != 0 {
		t.result, err = getTableLast(t.tableType, strconv.Itoa(lFlag), repFormat)
	} else if dFlag == "today" {
		t.result, err = getTableToday(t.tableType, repFormat)
	} else if dFlag == "current" {
		t.result, err = getTableCurrent(t.tableType, repFormat)
	} else if len(dFlag) == 10 {
		t.result, err = getTableDay(t.tableType, dFlag, repFormat)
	} else if len(dFlag) == 21 {
		t.result, err = getTableRange(t.tableType, dFlag, repFormat)
	}
	if err != nil {
		log.Fatal(err)
	}

	if repFormat != "xml" {
		if t.tableType != "C" {
			err = json.Unmarshal(t.result, &t.exchange)
		} else {
			err = json.Unmarshal(t.result, &t.exchangeC)
		}
		if err != nil {
			log.Fatal(err)
		}
	}

	return err
}

// getTableToday - function returns exchange rate table published today
// in JSON form, or error
func getTableToday(tableType string, repFormat string) ([]byte, error) {
	address := fmt.Sprintf("%s/tables/%s/today/?format=%s", baseAddressTable, tableType, repFormat)
	return getJSON(address)
}

// getTableCurrent - function returns current table of exchange rates
// (last published table) in JSON form, or error
func getTableCurrent(tableType string, repFormat string) ([]byte, error) {
	address := fmt.Sprintf("%s/tables/%s/?format=%s", baseAddressTable, tableType, repFormat)
	return getJSON(address)
}

// getTableDay - functions returns table of exchange rates
// on the given date (YYYY-MM-DD) in JSON form, or error
func getTableDay(tableType string, day string, repFormat string) ([]byte, error) {
	address := fmt.Sprintf("%s/tables/%s/%s/?format=%s", baseAddressTable, tableType, day, repFormat)
	return getJSON(address)
}

// getTableRange - function returns table of exchange rates  within
// the given date range (RRRR-MM-DD:RRRR-MM-DD) in JSON form, or error
func getTableRange(tableType string, day string, repFormat string) ([]byte, error) {
	var startDate string
	var stopDate string

	temp := strings.Split(day, ":")
	startDate = temp[0]
	stopDate = temp[1]

	address := fmt.Sprintf("%s/tables/%s/%s/%s/?format=%s", baseAddressTable, tableType, startDate, stopDate, repFormat)
	return getJSON(address)
}

// getTableLast - function returns last <last> tables of exchange rates
// in JSON form, or error
func getTableLast(tableType string, last string, repFormat string) ([]byte, error) {
	address := fmt.Sprintf("%s/tables/%s/last/%s/?format=%s", baseAddressTable, tableType, last, repFormat)
	return getJSON(address)
}

// GetPretty - function returns tables of exchange rates as
// formatted table,
// depending on the tableType field: for type A and B tables
// a column with an average rate is printed, for type C two columns:
// buy price and sell price
func (t *NBPTable) GetPretty() string {
	const padding = 3
	var builder strings.Builder
	var output string
	w := tabwriter.NewWriter(&builder, 0, 0, padding, ' ', tabwriter.Debug)

	if t.tableType != "C" {
		for _, item := range t.exchange {
			output += fmt.Sprintln()
			output += fmt.Sprintln(l.Get("Table type:")+"\t\t", item.Table)
			output += fmt.Sprintln(l.Get("Table number:")+"\t\t", item.No)
			output += fmt.Sprintln(l.Get("Publication date:")+"\t", item.EffectiveDate)
			output += fmt.Sprintln()

			fmt.Fprintln(w, l.Get("CODE \t NAME \t AVERAGE (PLN)"))
			fmt.Fprintln(w, l.Get("---- \t ---- \t -------------"))
			for _, currencyItem := range item.Rates {
				currencyValue := fmt.Sprintf("%.4f", currencyItem.Mid)
				fmt.Fprintln(w, currencyItem.Code+" \t "+currencyItem.Currency+" \t "+currencyValue)
			}
			w.Flush()
			output += builder.String()
			builder.Reset()
		}
	} else {
		for _, item := range t.exchangeC {
			output += fmt.Sprintln()
			output += fmt.Sprintln(l.Get("Table type:")+"\t\t", item.Table)
			output += fmt.Sprintln(l.Get("Table number:")+"\t\t", item.No)
			output += fmt.Sprintln(l.Get("Trading date:")+"\t\t", item.TradingDate)
			output += fmt.Sprintln(l.Get("Publication date:")+"\t", item.EffectiveDate)
			output += fmt.Sprintln()

			fmt.Fprintln(w, l.Get("CODE \t NAME \t BUY (PLN) \t SELL (PLN) "))
			fmt.Fprintln(w, l.Get("---- \t ---- \t --------- \t ---------- "))
			for _, currencyItem := range item.Rates {
				currencyValueBid := fmt.Sprintf("%.4f", currencyItem.Bid)
				currencyValueAsk := fmt.Sprintf("%.4f", currencyItem.Ask)
				fmt.Fprintln(w, currencyItem.Code+" \t "+currencyItem.Currency+" \t "+currencyValueBid+" \t "+currencyValueAsk)
			}
			w.Flush()
			output += builder.String()
			builder.Reset()
		}
	}

	return output
}

// GetCSV - function prints tables of exchange rates in the console,
// in the form of CSV (data separated by a comma), depending on the
// tableType field: for type A and B tables a column with an average
// rate is printed, for type C two columns: buy price and sell price
func (t *NBPTable) GetCSV() string {
	var tableNo string
	var output string = ""

	if t.tableType != "C" {
		output += fmt.Sprintln(l.Get("TABLE,CODE,NAME,AVERAGE (PLN)"))

		for _, item := range t.exchange {
			tableNo = item.No
			for _, currencyItem := range item.Rates {
				currencyValue := fmt.Sprintf("%.4f", currencyItem.Mid)
				output += fmt.Sprintln(tableNo + "," + currencyItem.Code + "," + currencyItem.Currency + "," + currencyValue)
			}
		}
	} else {
		output += fmt.Sprintln(l.Get("TABLE,CODE,NAME,BUY (PLN),SELL (PLN)"))

		for _, item := range t.exchangeC {
			tableNo = item.No
			for _, currencyItem := range item.Rates {
				currencyValueBid := fmt.Sprintf("%.4f", currencyItem.Bid)
				currencyValueAsk := fmt.Sprintf("%.4f", currencyItem.Ask)
				output += fmt.Sprintln(tableNo + "," + currencyItem.Code + "," + currencyItem.Currency + "," + currencyValueBid + "," + currencyValueAsk)
			}
			output += fmt.Sprintln()
		}
	}

	return output
}

// GetRaw - function returns just result of request (json or xml)
func (t *NBPTable) GetRaw() string {
	return string(t.result)
}
