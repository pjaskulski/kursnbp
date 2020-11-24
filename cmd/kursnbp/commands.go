// subcommands functions
package main

import (
	"fmt"
	"log"

	"github.com/pjaskulski/kursnbp/pkg/nbpapi"
)

// goldCommand - function for 'gold' command (prices of gold)
func goldCommand() {
	err := nbpapi.CheckArg("gold", cfg.tableFlag, cfg.dateFlag, cfg.lastFlag, cfg.outputFlag, cfg.codeFlag)
	if err != nil {
		log.Fatal(err)
	}

	nbpGold := nbpapi.NewGold()

	err = nbpGold.GetGold(cfg.dateFlag, cfg.lastFlag, cfg.repFormat)
	if err != nil {
		log.Fatal(err)
	}

	var output string

	switch cfg.outputFlag {
	case "table":
		output = nbpGold.GetPretty()
	case "json", "xml":
		output = nbpGold.GetRaw()
	case "csv":
		output = nbpGold.GetCSV()
	}

	fmt.Println(output)
}

// currencyCommand - function for 'currency' command (currency exchange rates)
func currencyCommand() {

	err := nbpapi.CheckArg("currency", cfg.tableFlag, cfg.dateFlag, cfg.lastFlag, cfg.outputFlag, cfg.codeFlag)
	if err != nil {
		log.Fatal(err)
	}

	nbpCurrency := nbpapi.NewCurrency(cfg.tableFlag)

	err = nbpCurrency.GetCurrency(cfg.dateFlag, cfg.lastFlag, cfg.codeFlag, cfg.repFormat)
	if err != nil {
		log.Fatal(err)
	}

	var output string

	switch cfg.outputFlag {
	case "table":
		output = nbpCurrency.GetPretty()
	case "json", "xml":
		output = nbpCurrency.GetRaw()
	case "csv":
		output = nbpCurrency.GetCSV()
	}

	fmt.Println(output)
}

// tableCommand - function for 'table' command (tables with exchange rates)
func tableCommand() {

	err := nbpapi.CheckArg("table", cfg.tableFlag, cfg.dateFlag, cfg.lastFlag, cfg.outputFlag, cfg.codeFlag)
	if err != nil {
		log.Fatal(err)
	}

	nbpTable := nbpapi.NewTable(cfg.tableFlag)

	err = nbpTable.GetTable(cfg.dateFlag, cfg.lastFlag, cfg.repFormat)
	if err != nil {
		log.Fatal(err)
	}

	var output string

	switch cfg.outputFlag {
	case "table":
		output = nbpTable.GetPretty()
	case "json", "xml":
		output = nbpTable.GetRaw()
	case "csv":
		output = nbpTable.GetCSV()
	}

	fmt.Println(output)
}
