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

// getTableAddress - build download address depending on previously
// verified input parameters (--table, --date or --last)
func getTableAddress(tableType string, dFlag string, lFlag int) string {
	var address string

	if lFlag != 0 {
		address = queryTableLast(tableType, strconv.Itoa(lFlag))
	} else if dFlag == "today" {
		address = queryTableToday(tableType)
	} else if dFlag == "current" {
		address = queryTableCurrent(tableType)
	} else if len(dFlag) == 10 {
		address = queryTableDay(tableType, dFlag)
	} else if len(dFlag) == 21 {
		address = queryTableRange(tableType, dFlag)
	}

	return address
}

// GetTableRaw - function downloads data in json or xml form
func (t *NBPTable) GetTableRaw(dFlag string, lFlag int, repFormat string) error {
	var err error

	address := getTableAddress(t.tableType, dFlag, lFlag)
	t.result, err = getData(address, repFormat)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

// GetTableByDate - function downloads and writes data to exchange (exchangeC) slice,
// raw data (json) still available in result field
func (t *NBPTable) GetTableByDate(dFlag string) error {
	var err error

	address := getTableAddress(t.tableType, dFlag, 0)
	t.result, err = getData(address, "json")
	if err != nil {
		log.Fatal(err)
	}

	if t.tableType != "C" {
		err = json.Unmarshal(t.result, &t.exchange)
	} else {
		err = json.Unmarshal(t.result, &t.exchangeC)
	}
	if err != nil {
		log.Fatal(err)
	}

	return err
}

// GetTableLast - function downloads and writes data to exchange (exchangeC) slice,
// raw data (json) still available in result field
func (t *NBPTable) GetTableLast(lFlag int) error {
	var err error

	address := getTableAddress(t.tableType, "", lFlag)
	t.result, err = getData(address, "json")
	if err != nil {
		log.Fatal(err)
	}

	if t.tableType != "C" {
		err = json.Unmarshal(t.result, &t.exchange)
	} else {
		err = json.Unmarshal(t.result, &t.exchangeC)
	}
	if err != nil {
		log.Fatal(err)
	}

	return err
}

// queryTableToday - returns query: exchange rate table published today
func queryTableToday(tableType string) string {
	return fmt.Sprintf("%s/tables/%s/today/", baseAddressTable, tableType)
}

// queryTableCurrent - returns query: current table of exchange rates
// (last published table)
func queryTableCurrent(tableType string) string {
	return fmt.Sprintf("%s/tables/%s/", baseAddressTable, tableType)
}

// queryTableDay - returns query: table of exchange rates
// on the given date (YYYY-MM-DD)
func queryTableDay(tableType string, day string) string {
	return fmt.Sprintf("%s/tables/%s/%s/", baseAddressTable, tableType, day)
}

// queryTableRange - returns query: table of exchange rates  within
// the given date range (RRRR-MM-DD:RRRR-MM-DD)
func queryTableRange(tableType string, day string) string {
	var startDate string
	var stopDate string

	temp := strings.Split(day, ":")
	startDate = temp[0]
	stopDate = temp[1]

	address := fmt.Sprintf("%s/tables/%s/%s/%s/", baseAddressTable, tableType, startDate, stopDate)
	return address
}

// queryTableLast - returns query: last <number> tables of exchange rates
func queryTableLast(tableType string, last string) string {
	return fmt.Sprintf("%s/tables/%s/last/%s/", baseAddressTable, tableType, last)
}

// GetPrettyOutput - function returns tables of exchange rates as
// formatted table, depending on the tableType field: for type A and B tables
// a column with an average rate is printed, for type C two columns:
// buy price and sell price
func (t *NBPTable) GetPrettyOutput() string {
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

// GetCSVOutput - function prints tables of exchange rates in the console,
// in the form of CSV (data separated by a comma), depending on the
// tableType field: for type A and B tables a column with an average
// rate is printed, for type C two columns: buy price and sell price
func (t *NBPTable) GetCSVOutput() string {
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
		}
	}

	return output
}

// GetRawOutput - function returns just result of request (json or xml)
func (t *NBPTable) GetRawOutput() string {
	return string(t.result)
}
