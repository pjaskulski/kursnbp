// subcommands functions
package main

import (
	"fmt"
	"log"
)

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

// currencyCommand - function for 'currency' command (currency exchange rates)
func currencyCommand() {
	var result []byte

	err := checkArg("currency", tableFlag, dateFlag, lastFlag, outputFlag, codeFlag)
	if err != nil {
		log.Fatal(err)
	}

	result, err = getCurrency(tableFlag, dateFlag, lastFlag, codeFlag)
	if err != nil {
		log.Fatal(err)
	}

	switch outputFlag {
	case "table":
		printCurrency(result, tableFlag)
	case "json", "xml":
		fmt.Println(string(result))
	case "csv":
		printCurrencyCSV(result, tableFlag)
	}
}

// tableCommand - function for 'table' command (tables with exchange rates)
func tableCommand() {
	var result []byte

	err := checkArg("table", tableFlag, dateFlag, lastFlag, outputFlag, codeFlag)
	if err != nil {
		log.Fatal(err)
	}

	result, err = getTable(tableFlag, dateFlag, lastFlag)
	if err != nil {
		log.Fatal(err)
	}

	switch outputFlag {
	case "table":
		printTable(result, tableFlag)
	case "json", "xml":
		fmt.Println(string(result))
	case "csv":
		printTableCSV(result, tableFlag)
	}
}
