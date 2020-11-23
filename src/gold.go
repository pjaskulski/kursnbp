// gold' subcommand support - downloading gold prices in the JSON format

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

type rateGold struct {
	Data string  `json:"data"`
	Cena float64 `json:"cena"`
}

// goldCommand - function for 'gold' command (prices of gold)
func goldCommand() {
	var result []byte

	err := checkArg("gold", tableFlag, dateFlag, lastFlag, outputFlag, codeFlag)
	if err != nil {
		log.Fatal(err)
	}

	result, err = getGold(dateFlag, lastFlag)
	if err != nil {
		log.Fatal(err)
	}

	switch outputFlag {
	case "table":
		printGold(result)
	case "json", "xml":
		fmt.Println(string(result))
	case "csv":
		printGoldCSV(result)
	}
}

// getGold - main function for gold prices, selects
// a data download variant depending on previously
// verified input parameters (--date or --last)
func getGold(dFlag string, lFlag int) ([]byte, error) {
	var result []byte
	var err error

	if lFlag != 0 {
		result, err = getGoldLast(strconv.Itoa(lFlag))
	} else if dFlag == "today" {
		result, err = getGoldToday()
	} else if dFlag == "current" {
		result, err = getGoldCurrent()
	} else if len(dFlag) == 10 {
		result, err = getGoldDay(dFlag)
	} else if len(dFlag) == 21 {
		result, err = getGoldRange(dFlag)
	}

	return result, err
}

// getGoldToday - function returns today's gold price
// in json form, or error
func getGoldToday() ([]byte, error) {
	address := fmt.Sprintf(baseAddressGold + "/today?format=" + repFormat)
	return getJSON(address)
}

// getGoldCurrent - function returns current gold price
// (last published price) in json form, or error
func getGoldCurrent() ([]byte, error) {
	address := fmt.Sprintf(baseAddressGold + "?format=" + repFormat)
	return getJSON(address)
}

// getGoldLast - function returns last <last> gold prices
// in json form, or error
func getGoldLast(last string) ([]byte, error) {
	address := fmt.Sprintf(baseAddressGold+"/last/%s?format="+repFormat, last)
	return getJSON(address)
}

// getGoldDay - function returns gold price on the given date (RRRR-MM-DD)
// in json form, or error
func getGoldDay(day string) ([]byte, error) {
	address := fmt.Sprintf(baseAddressGold+"/%s?format="+repFormat, day)
	return getJSON(address)
}

// getGoldRange - function returns gold prices within the given date range
// (RRRR-MM-DD:RRRR-MM-DD) in json form, or error
func getGoldRange(day string) ([]byte, error) {
	var startDate string
	var stopDate string

	temp := strings.Split(day, ":")
	if len(temp) != 2 {
		log.Fatal(errors.New(l.Get("Invalid date range format")))
	}

	startDate = temp[0]
	stopDate = temp[1]

	address := fmt.Sprintf(baseAddressGold+"/%s/%s?format="+repFormat, startDate, stopDate)
	return getJSON(address)
}

// printGold - functions displays a formatted table of gold prices
// in the console window
func printGold(result []byte) {
	var nbpGold []rateGold
	err := json.Unmarshal(result, &nbpGold)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(appName, "-", appDesc)

	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.Debug)

	fmt.Println()
	fmt.Println(l.Get("The price of 1g of gold (of 1000 millesimal fineness)"))
	fmt.Println()

	fmt.Fprintln(w, l.Get("DATE \t PRICE (PLN)"))
	fmt.Fprintln(w, l.Get("---- \t ----- "))
	for _, goldItem := range nbpGold {
		goldValue := fmt.Sprintf("%.4f", goldItem.Cena)
		fmt.Fprintln(w, goldItem.Data+" \t "+goldValue)
	}
	w.Flush()

	fmt.Println()
}

// printGoldCSV - function prints gold prices in CSV format
// (comma separated data)
func printGoldCSV(result []byte) {
	var nbpGold []rateGold
	err := json.Unmarshal(result, &nbpGold)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(l.Get("DATE,PRICE (PLN)"))
	for _, goldItem := range nbpGold {
		goldValue := fmt.Sprintf("%.4f", goldItem.Cena)
		fmt.Println(goldItem.Data + "," + goldValue)
	}

	fmt.Println()
}
