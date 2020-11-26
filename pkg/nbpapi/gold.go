// gold' subcommand support - downloading gold prices in the JSON format

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
	baseAddressGold string = "http://api.nbp.pl/api/cenyzlota"
)

type rateGold struct {
	Data  string  `json:"data"`
	Price float64 `json:"cena"`
}

// NBPGold type
type NBPGold struct {
	GoldRates []rateGold
	result    []byte
}

// NewGold - function creates new gold type
func NewGold() *NBPGold {
	return &NBPGold{}
}

// getGoldAddress - build download address depending on previously
// verified input parameters (--date or --last)
func getGoldAddress(dFlag string, lFlag int) string {
	var address string

	if lFlag != 0 {
		address = queryGoldLast(strconv.Itoa(lFlag))
	} else if dFlag == "today" {
		address = queryGoldToday()
	} else if dFlag == "current" {
		address = queryGoldCurrent()
	} else if len(dFlag) == 10 {
		address = queryGoldDay(dFlag)
	} else if len(dFlag) == 21 {
		address = queryGoldRange(dFlag)
	}

	return address
}

// GetGoldRaw - function downloads data in json or xml form
func (g *NBPGold) GetGoldRaw(dFlag string, lFlag int, repFormat string) error {
	var err error

	address := getGoldAddress(dFlag, lFlag)
	g.result, err = getData(address, repFormat)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

// GetGoldByDate - function downloads and writes data to goldRates slice,
// raw data (json) still available in result field
func (g *NBPGold) GetGoldByDate(dFlag string) error {
	var err error

	address := getGoldAddress(dFlag, 0)
	g.result, err = getData(address, "json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(g.result, &g.GoldRates)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

// GetGoldLast - function downloads and writes data to goldRates slice,
// raw data (json) still available in result field
func (g *NBPGold) GetGoldLast(lFlag int) error {
	var err error

	address := getGoldAddress("", lFlag)
	g.result, err = getData(address, "json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(g.result, &g.GoldRates)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

// queryGoldToday - returns query: today's gold price
func queryGoldToday() string {
	return fmt.Sprintf("%s/today", baseAddressGold)
}

// queryGoldCurrent - returns query: current gold price
// (last published price)
func queryGoldCurrent() string {
	return baseAddressGold
}

// queryGoldLast - returns query: last <number> gold prices
func queryGoldLast(last string) string {
	return fmt.Sprintf("%s/last/%s", baseAddressGold, last)
}

// queryGoldDay - function returns gold price on the given date (RRRR-MM-DD)
// in json/xml form, or error
func queryGoldDay(day string) string {
	return fmt.Sprintf("%s/%s", baseAddressGold, day)
}

// queryGoldRange - returns query: gold prices within the given date range
// (RRRR-MM-DD:RRRR-MM-DD) in json/xml form, or error
func queryGoldRange(day string) string {
	var startDate string
	var stopDate string

	temp := strings.Split(day, ":")
	startDate = temp[0]
	stopDate = temp[1]

	address := fmt.Sprintf("%s/%s/%s", baseAddressGold, startDate, stopDate)
	return address
}

// GetPrettyOutput - function returns a formatted table of gold prices
func (g *NBPGold) GetPrettyOutput() string {

	const padding = 3
	var builder strings.Builder
	w := tabwriter.NewWriter(&builder, 0, 0, padding, ' ', tabwriter.Debug)

	fmt.Fprintln(w)
	fmt.Fprintln(w, l.Get("The price of 1g of gold (of 1000 millesimal fineness)"))
	fmt.Fprintln(w)

	fmt.Fprintln(w, l.Get("DATE \t PRICE (PLN)"))
	fmt.Fprintln(w, l.Get("---- \t ----------- "))
	for _, goldItem := range g.GoldRates {
		goldValue := fmt.Sprintf("%.4f", goldItem.Price)
		fmt.Fprintln(w, goldItem.Data+" \t "+goldValue)
	}
	w.Flush()

	return builder.String()
}

// GetCSVOutput - function returns prices of gold in CSV format
// (comma separated data)
func (g *NBPGold) GetCSVOutput() string {
	var output string = ""

	output += fmt.Sprintln(l.Get("DATE,PRICE (PLN)"))
	for _, goldItem := range g.GoldRates {
		goldValue := fmt.Sprintf("%.4f", goldItem.Price)
		output += fmt.Sprintln(goldItem.Data + "," + goldValue)
	}

	return output
}

// GetRawOutput - function returns just result of request (json or xml)
func (g *NBPGold) GetRawOutput() string {
	return string(g.result)
}
