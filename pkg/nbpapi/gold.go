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
	Data string  `json:"data"`
	Cena float64 `json:"cena"`
}

// NBPGold type
type NBPGold struct {
	goldRates []rateGold
	result    []byte
}

// NewGold - function creates new gold type
func NewGold() *NBPGold {
	return &NBPGold{}
}

// GetGold - main function for gold prices, selects
// a data download variant depending on previously
// verified input parameters (--date or --last)
func (g *NBPGold) GetGold(dFlag string, lFlag int, repFormat string) error {
	var err error

	if lFlag != 0 {
		g.result, err = getGoldLast(strconv.Itoa(lFlag), repFormat)
	} else if dFlag == "today" {
		g.result, err = getGoldToday(repFormat)
	} else if dFlag == "current" {
		g.result, err = getGoldCurrent(repFormat)
	} else if len(dFlag) == 10 {
		g.result, err = getGoldDay(dFlag, repFormat)
	} else if len(dFlag) == 21 {
		g.result, err = getGoldRange(dFlag, repFormat)
	}

	if err != nil {
		log.Fatal(err)
	}

	if repFormat != "xml" {
		err = json.Unmarshal(g.result, &g.goldRates)
		if err != nil {
			log.Fatal(err)
		}
	}

	return err
}

// getGoldToday - function returns today's gold price
// in json form, or error
func getGoldToday(repFormat string) ([]byte, error) {
	address := fmt.Sprintf("%s/today?format=%s", baseAddressGold, repFormat)
	return getJSON(address)
}

// getGoldCurrent - function returns current gold price
// (last published price) in json form, or error
func getGoldCurrent(repFormat string) ([]byte, error) {
	address := fmt.Sprintf("%s?format=%s", baseAddressGold, repFormat)
	return getJSON(address)
}

// getGoldLast - function returns last <last> gold prices
// in json form, or error
func getGoldLast(last string, repFormat string) ([]byte, error) {
	address := fmt.Sprintf("%s/last/%s?format=%s", baseAddressGold, last, repFormat)
	return getJSON(address)
}

// getGoldDay - function returns gold price on the given date (RRRR-MM-DD)
// in json form, or error
func getGoldDay(day string, repFormat string) ([]byte, error) {
	address := fmt.Sprintf("%s/%s?format=%s", baseAddressGold, day, repFormat)
	return getJSON(address)
}

// getGoldRange - function returns gold prices within the given date range
// (RRRR-MM-DD:RRRR-MM-DD) in json form, or error
func getGoldRange(day string, repFormat string) ([]byte, error) {
	var startDate string
	var stopDate string

	temp := strings.Split(day, ":")
	startDate = temp[0]
	stopDate = temp[1]

	address := fmt.Sprintf("%s/%s/%s?format=%s", baseAddressGold, startDate, stopDate, repFormat)
	return getJSON(address)
}

// GetPretty - function returns a formatted table of gold prices
func (g *NBPGold) GetPretty() string {

	const padding = 3
	var builder strings.Builder
	w := tabwriter.NewWriter(&builder, 0, 0, padding, ' ', tabwriter.Debug)

	fmt.Fprintln(w)
	fmt.Fprintln(w, l.Get("The price of 1g of gold (of 1000 millesimal fineness)"))
	fmt.Fprintln(w)

	fmt.Fprintln(w, l.Get("DATE \t PRICE (PLN)"))
	fmt.Fprintln(w, l.Get("---- \t ----------- "))
	for _, goldItem := range g.goldRates {
		goldValue := fmt.Sprintf("%.4f", goldItem.Cena)
		fmt.Fprintln(w, goldItem.Data+" \t "+goldValue)
	}
	w.Flush()

	return builder.String()
}

// GetCSV - function returns prices of gold in CSV format
// (comma separated data)
func (g *NBPGold) GetCSV() string {
	var output string = ""

	output += fmt.Sprintln(l.Get("DATE,PRICE (PLN)"))
	for _, goldItem := range g.goldRates {
		goldValue := fmt.Sprintf("%.4f", goldItem.Cena)
		output += fmt.Sprintln(goldItem.Data + "," + goldValue)
	}

	return output
}

// GetRaw - function returns just result of request (json or xml)
func (g *NBPGold) GetRaw() string {
	return string(g.result)
}
