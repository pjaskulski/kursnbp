package main

import (
	"os"
	"strings"

	"github.com/integrii/flaggy"
	"github.com/pjaskulski/kursnbp/pkg/nbpapi"
)

// app name, version and description
var (
	version string = "0.3.1"
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
var cfg struct {
	tableFlag  string
	dateFlag   string
	outputFlag string
	lastFlag   int
	codeFlag   string
	langFlag   string
	repFormat  string
}

func init() {
	// ini cfg (default values)
	cfg.tableFlag = "A"
	cfg.outputFlag = "table"
	cfg.langFlag = "en"
	cfg.repFormat = "json"

	// command line support through the flaggy package
	flaggy.SetName(appName)
	flaggy.SetDescription(appDesc)

	flaggy.DefaultParser.ShowHelpOnUnexpected = true
	flaggy.DefaultParser.AdditionalHelpPrepend = "https://pjaskulski.github.io/kursnbp/"

	// table subcommand
	cmdTable = flaggy.NewSubcommand("table")
	cmdTable.Description = "returns a table of exchange rates (or a series of tables)"
	cmdTable.String(&cfg.tableFlag, "t", "table", "type of exchange rate table, 'A', 'B' or 'C'")
	cmdTable.String(&cfg.dateFlag, "d", "date", "date 'YYYYY-MM-DD', or range of dates 'YYYY-MM-DD:YYYY-MM-DD', or 'today' or 'current' (default: current)")
	cmdTable.Int(&cfg.lastFlag, "l", "last", "a series of the last <number> of exchange rate tables")
	cmdTable.String(&cfg.outputFlag, "o", "output", "output format: 'table', 'json', 'csv'")
	cmdTable.String(&cfg.langFlag, "i", "lang", "output language: 'en', 'pl'")
	flaggy.AttachSubcommand(cmdTable, 1)

	// currency subcommand
	cmdCurrency = flaggy.NewSubcommand("currency")
	cmdCurrency.Description = "returns the rate of the indicated currency or a series of rates"
	cmdCurrency.String(&cfg.tableFlag, "t", "table", "type of exchange rate table, 'A', 'B' or 'C'")
	cmdCurrency.String(&cfg.dateFlag, "d", "date", "date 'YYYYY-MM-DD', or range of dates 'YYYY-MM-DD:YYYY-MM-DD', or 'today' or 'current' (default: current)")
	cmdCurrency.String(&cfg.codeFlag, "c", "code", "currency code according to ISO 4217 e.g. CHF")
	cmdCurrency.Int(&cfg.lastFlag, "l", "last", "series of last <number> exchange rates of the indicated currency")
	cmdCurrency.String(&cfg.outputFlag, "o", "output", "output format: 'table', 'json', 'csv'")
	cmdCurrency.String(&cfg.langFlag, "i", "lang", "output language: 'en', 'pl'")
	flaggy.AttachSubcommand(cmdCurrency, 1)

	// gold subcommand
	cmdGold = flaggy.NewSubcommand("gold")
	cmdGold.Description = "returns a gold price or a series of gold price quotations"
	cmdGold.String(&cfg.dateFlag, "d", "date", "date 'YYYYY-MM-DD', or range of dates 'YYYY-MM-DD:YYYY-MM-DD', or 'today' or 'current' (default: current)")
	cmdGold.Int(&cfg.lastFlag, "l", "last", "a series of recent <number> gold price quotations")
	cmdGold.String(&cfg.outputFlag, "o", "output", "output format: 'table', 'json', 'csv'")
	cmdGold.String(&cfg.langFlag, "i", "lang", "output language: 'en', 'pl'")
	flaggy.AttachSubcommand(cmdGold, 1)

	flaggy.SetVersion(version)
	flaggy.Parse()

	// modifications to received as command line parameters flag values, the default
	// value of the flag --date is set this way, because the flags --date and --last
	// are alternative, the default value --date makes sense only if the user has not
	// set --last
	if cfg.lastFlag == 0 && cfg.dateFlag == "" {
		cfg.dateFlag = "current"
	}
	// modifications to the flag values: the characters of the --table and --code values
	// are changed to upper, therefore it is acceptable to call --code=chf, or --code=CHf,
	// the application will support such call correctly
	if cfg.tableFlag != "" {
		cfg.tableFlag = strings.ToUpper(cfg.tableFlag)
	}
	if cfg.codeFlag != "" {
		cfg.codeFlag = strings.ToUpper(cfg.codeFlag)
	}
	// modifications to the flag values: the characters of the --lang values
	// are changed to lower, therefore it is acceptable to call --lang=PL, or --lang=Pl,
	// the application will support such call correctly
	if cfg.langFlag != "" {
		cfg.langFlag = strings.ToLower(cfg.langFlag)
	}

	if cfg.outputFlag == "xml" {
		cfg.repFormat = "xml"
	}
}

// kursnbp - command line tool for downloading exchange rates and gold prices
// from the website of the National Bank of Poland (http://api.nbp.pl/en.html)
func main() {

	// set output message language based on --lang flag, English is default,
	// if flag --lang is different than 'pl' or 'en' English is set
	nbpapi.SetLang(cfg.langFlag)

	if cmdTable.Used {
		tableCommand()
	} else if cmdCurrency.Used {
		currencyCommand()
	} else if cmdGold.Used {
		goldCommand()
	} else {
		// if no correct subcommand is given, a general help is displayed
		// and the program ends
		flaggy.ShowHelp("")
		os.Exit(1)
	}
}
