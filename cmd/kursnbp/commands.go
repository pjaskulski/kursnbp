// subcommands functions
package main

import (
	"log"
)

// goldCommand - function for 'gold' command (prices of gold)
func goldCommand() {
	err := checkArg("gold", tableFlag, dateFlag, lastFlag, outputFlag, codeFlag)
	if err != nil {
		log.Fatal(err)
	}

	var gold Gold

	err = gold.GetGold(dateFlag, lastFlag)
	if err != nil {
		log.Fatal(err)
	}

	switch outputFlag {
	case "table":
		gold.PrintGold()
	case "json", "xml":
		gold.PrintResult()
	case "csv":
		gold.PrintGoldCSV()
	}
}

// currencyCommand - function for 'currency' command (currency exchange rates)
func currencyCommand() {

	err := checkArg("currency", tableFlag, dateFlag, lastFlag, outputFlag, codeFlag)
	if err != nil {
		log.Fatal(err)
	}

	var nbpCurrency NBPCurrency
	nbpCurrency.tableType = tableFlag

	err = nbpCurrency.GetCurrency(dateFlag, lastFlag, codeFlag)
	if err != nil {
		log.Fatal(err)
	}

	switch outputFlag {
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

	err := checkArg("table", tableFlag, dateFlag, lastFlag, outputFlag, codeFlag)
	if err != nil {
		log.Fatal(err)
	}

	var nbpTable NBPTable
	nbpTable.tableType = tableFlag

	err = nbpTable.GetTable(dateFlag, lastFlag)
	if err != nil {
		log.Fatal(err)
	}

	switch outputFlag {
	case "table":
		nbpTable.PrintTable()
	case "json", "xml":
		nbpTable.PrintResult()
	case "csv":
		nbpTable.PrintTableCSV()
	}
}
