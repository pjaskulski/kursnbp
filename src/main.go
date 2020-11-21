package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/integrii/flaggy"
)

// base addresses of the NBP API service for currency rates and gold prices
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

// subcommands
var (
	cmdTable    *flaggy.Subcommand
	cmdCurrency *flaggy.Subcommand
	cmdGold     *flaggy.Subcommand
)

// flags
var (
	tableFlag  string = "A"
	dateFlag   string = ""
	outputFlag string = "table"
	lastFlag   int
	codeFlag   string = ""
	langFlag   string = "en"
)

func init() {
	// command line support through the flaggy package
	flaggy.SetName(appName)
	flaggy.SetDescription(appDesc)

	flaggy.DefaultParser.ShowHelpOnUnexpected = true
	flaggy.DefaultParser.AdditionalHelpPrepend = "http://github.com/pjaskulski/kursnbp"

	// table subcommand
	cmdTable = flaggy.NewSubcommand("table")
	cmdTable.Description = "returns a table of exchange rates (or a series of tables)"
	cmdTable.String(&tableFlag, "t", "table", "type of exchange rate table, 'A', 'B' or 'C'")
	cmdTable.String(&dateFlag, "d", "date", "date 'YYYYY-MM-DD', or date range 'YYYY-MM-DD:YYYY-MM-DD', or 'today' or 'current' (default: current)")
	cmdTable.Int(&lastFlag, "l", "last", "a series of the last <number> of exchange rate tables")
	cmdTable.String(&outputFlag, "o", "output", "output format: 'table', 'json', 'csv'")
	cmdTable.String(&langFlag, "i", "lang", "output language: 'en', 'pl'")
	flaggy.AttachSubcommand(cmdTable, 1)

	// currency subcommand
	cmdCurrency = flaggy.NewSubcommand("currency")
	cmdCurrency.Description = "returns the rate of the indicated currency or a series of rates"
	cmdCurrency.String(&tableFlag, "t", "table", "type of exchange rate table, 'A', 'B' or 'C'")
	cmdCurrency.String(&dateFlag, "d", "date", "date 'YYYYY-MM-DD', or date range 'YYYY-MM-DD:YYYY-MM-DD', or 'today' or 'current' (default: current)")
	cmdCurrency.String(&codeFlag, "c", "code", "currency code according to ISO 4217 e.g. CHF")
	cmdCurrency.Int(&lastFlag, "l", "last", "series of last <number> exchange rates of the indicated currency")
	cmdCurrency.String(&outputFlag, "o", "output", "output format: 'table', 'json', 'csv'")
	cmdCurrency.String(&langFlag, "i", "lang", "output language: 'en', 'pl'")
	flaggy.AttachSubcommand(cmdCurrency, 1)

	// gold subcommand
	cmdGold = flaggy.NewSubcommand("gold")
	cmdGold.Description = "returns a gold price or a series of gold price quotations"
	cmdGold.String(&dateFlag, "d", "date", "date 'YYYYY-MM-DD', or date range 'YYYY-MM-DD:YYYY-MM-DD', or 'today' or 'current' (default: current)")
	cmdGold.Int(&lastFlag, "l", "last", "a series of recent <number> gold price quotations")
	cmdGold.String(&outputFlag, "o", "output", "output format: 'table', 'json', 'csv'")
	cmdGold.String(&langFlag, "i", "lang", "output language: 'en', 'pl'")
	flaggy.AttachSubcommand(cmdGold, 1)

	flaggy.SetVersion(version)
	flaggy.Parse()

	// modifications to received as command line parameters flag values, the default
	// value of the flag --date is set this way, because the flags --date and --last
	// are alternative, the default value --date makes sense only if the user has not
	// set --last
	if lastFlag == 0 && dateFlag == "" {
		dateFlag = "current"
	}
	// modifications to the flag values: the characters of the --table and --code values
	// are changed to upper, therefore it is acceptable to call --code=chf, or --code=CHf,
	// the application will support such call correctly
	if tableFlag != "" {
		tableFlag = strings.ToUpper(tableFlag)
	}
	if codeFlag != "" {
		codeFlag = strings.ToUpper(codeFlag)
	}
	// modifications to the flag values: the characters of the --lang values
	// are changed to lower, therefore it is acceptable to call --lang=PL, or --lang=Pl,
	// the application will support such call correctly
	if langFlag != "" {
		langFlag = strings.ToLower(langFlag)
	}
}

// kursnbp - command line tool for downloading exchange rates and gold prices
// from the website of the National Bank of Poland (http://api.nbp.pl/en.html)
func main() {
	var result []byte
	var err error

	// set output message language
	if langFlag == "pl" {
		l = langTexts["pl"]
	} else {
		l = langTexts["en"]
	}

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
		// if no correct subcommand is given, a general help is displayed
		// and the program ends
		flaggy.ShowHelp("")
		os.Exit(1)
	}
}
