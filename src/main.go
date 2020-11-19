package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/integrii/flaggy"
)

// adresy bazowe dla kursów walut i cen złota
const (
	baseAddress     string = "http://api.nbp.pl/api/exchangerates"
	baseAddressGold string = "http://api.nbp.pl/api/cenyzlota"
)

// app name, version and description
var (
	version string = "0.2.0"
	appName string = "kursNBP"
	appDesc string = "tool for downloading exchange rates and gold prices from the website of the National Bank of Poland"
)

// polecenia
var (
	cmdTable    *flaggy.Subcommand
	cmdCurrency *flaggy.Subcommand
	cmdGold     *flaggy.Subcommand
)

// flagi
var (
	tableFlag  string = "A"
	dateFlag   string = ""
	outputFlag string = "table"
	lastFlag   int
	codeFlag   string = ""
)

func init() {
	flaggy.SetName(appName)
	flaggy.SetDescription(appDesc)

	flaggy.DefaultParser.ShowHelpOnUnexpected = true
	flaggy.DefaultParser.AdditionalHelpPrepend = "http://github.com/pjaskulski/kursnbp"

	// table subcommand
	cmdTable = flaggy.NewSubcommand("table")
	cmdTable.Description = "zwraca tabelę kursów wymiany walut (lub serię tabel)"
	cmdTable.String(&tableFlag, "t", "table", "typ tabeli kursów, 'A', 'B' lub 'C'")
	cmdTable.String(&dateFlag, "d", "date", "data 'RRRR-MM-DD', lub zakres dat 'RRRR-MM-DD:RRRR-MM-DD', lub 'today' lub 'current' (defalut: current)")
	cmdTable.Int(&lastFlag, "l", "last", "seria <number> ostatnich tabel kursów")
	cmdTable.String(&outputFlag, "o", "output", "format wyjścia: 'table', 'json', 'csv'")

	flaggy.AttachSubcommand(cmdTable, 1)

	// currency subcommand
	cmdCurrency = flaggy.NewSubcommand("currency")
	cmdCurrency.Description = "zwraca kurs wskazanej waluty lub serię kursów"
	cmdCurrency.String(&tableFlag, "t", "table", "typ tabeli kursów, A, B lub C")
	cmdCurrency.String(&dateFlag, "d", "date", "data RRRR-MM-DD, lub zakres dat RRRR-MM-DD:RRRR-MM-DD, lub 'today' lub 'current' (defalut: current)")
	cmdCurrency.String(&codeFlag, "c", "code", "kod waluty zgodny z ISO 4217 np. CHF")
	cmdCurrency.Int(&lastFlag, "l", "last", "seria <number> ostatnich kursów waluty")
	cmdCurrency.String(&outputFlag, "o", "output", "format wyjścia: table, json, csv")
	flaggy.AttachSubcommand(cmdCurrency, 1)

	// gold subcommand
	cmdGold = flaggy.NewSubcommand("gold")
	cmdGold.Description = "zwraca cenę złota lub serię notowań cen złota"
	cmdGold.String(&dateFlag, "d", "date", "data 'RRRR-MM-DD', lub zakres dat 'RRRR-MM-DD:RRRR-MM-DD', lub 'today' lub 'current' (defalut: current)")
	cmdGold.Int(&lastFlag, "l", "last", "seria <number> ostatnich notowań cen złota")
	cmdGold.String(&outputFlag, "o", "output", "format wyjścia: table, json, csv")
	flaggy.AttachSubcommand(cmdGold, 1)

	flaggy.SetVersion(version)
	flaggy.Parse()

	if lastFlag == 0 && dateFlag == "" {
		dateFlag = "current"
	}
	if tableFlag != "" {
		tableFlag = strings.ToUpper(tableFlag)
	}
	if codeFlag != "" {
		codeFlag = strings.ToUpper(codeFlag)
	}
}

// kursnbp - command line tool for downloading exchange rates and gold prices
// from the website of the National Bank of Poland (http://api.nbp.pl/en.html)
func main() {
	var result []byte
	var err error

	if cmdTable.Used {
		err = checkArg("table", tableFlag, dateFlag, lastFlag, outputFlag, codeFlag)
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
		case "json":
			fmt.Println(string(result))
		case "csv":
			printTableCSV(result, tableFlag)
		}

	} else if cmdCurrency.Used {
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
		case "json":
			fmt.Println(string(result))
		case "csv":
			printCurrencyCSV(result, tableFlag)
		}

	} else if cmdGold.Used {
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
		case "json":
			fmt.Println(string(result))
		case "csv":
			printGoldCSV(result)
		}

	} else {
		flaggy.ShowHelp("")
		os.Exit(1)
	}
}
