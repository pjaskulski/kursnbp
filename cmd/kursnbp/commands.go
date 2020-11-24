// subcommands functions
package main

import (
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

	switch cfg.outputFlag {
	case "table":
		nbpGold.PrintGold()
	case "json", "xml":
		nbpGold.PrintResult()
	case "csv":
		nbpGold.PrintGoldCSV()
	}
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

	switch cfg.outputFlag {
	case "table":
		nbpCurrency.PrintCurrency()
	case "json", "xml":
		nbpCurrency.PrintResult()
	case "csv":
		nbpCurrency.PrintCurrencyCSV()
	}
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

	switch cfg.outputFlag {
	case "table":
		nbpTable.PrintTable()
	case "json", "xml":
		nbpTable.PrintResult()
	case "csv":
		nbpTable.PrintTableCSV()
	}
}
