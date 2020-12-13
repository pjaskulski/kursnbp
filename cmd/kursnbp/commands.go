// subcommands functions
package main

import (
	"fmt"
	"log"

	"github.com/atotto/clipboard"
	"github.com/pjaskulski/nbpapi"
)

// goldCommand - function for 'gold' command (prices of gold)
func goldCommand() {
	var err error

	err = checkArg("gold", cfg.tableFlag, cfg.dateFlag, cfg.lastFlag, cfg.outputFlag, cfg.codeFlag)
	if err != nil {
		log.Fatal(err)
	}

	nbpGold := nbpapi.NewGold()

	if cfg.outputFlag == "xml" || cfg.outputFlag == "json" {
		err = nbpGold.GoldRaw(cfg.dateFlag, cfg.lastFlag, cfg.repFormat)
	} else if cfg.lastFlag > 0 {
		err = nbpGold.GoldLast(cfg.lastFlag)
	} else {
		err = nbpGold.GoldByDate(cfg.dateFlag)
	}
	if err != nil {
		log.Fatal(err)
	}

	var output string

	switch cfg.outputFlag {
	case "table":
		output = nbpGold.CreatePrettyOutput(cfg.langFlag)
	case "json", "xml":
		output = nbpGold.CreateRawOutput()
	case "csv":
		output = nbpGold.CreateCSVOutput(cfg.langFlag)
	}

	if cfg.clipFlag {
		err = clipboard.WriteAll(output)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println(output)
	}
}

// currencyCommand - function for 'currency' command (currency exchange rates)
func currencyCommand() {
	var err error

	err = checkArg("currency", cfg.tableFlag, cfg.dateFlag, cfg.lastFlag, cfg.outputFlag, cfg.codeFlag)
	if err != nil {
		log.Fatal(err)
	}

	nbpCurrency := nbpapi.NewCurrency(cfg.tableFlag)

	if cfg.outputFlag == "xml" || cfg.outputFlag == "json" {
		err = nbpCurrency.CurrencyRaw(cfg.codeFlag, cfg.dateFlag, cfg.lastFlag, cfg.repFormat)
	} else if cfg.lastFlag > 0 {
		err = nbpCurrency.CurrencyLast(cfg.codeFlag, cfg.lastFlag)
	} else {
		err = nbpCurrency.CurrencyByDate(cfg.codeFlag, cfg.dateFlag)
	}
	if err != nil {
		log.Fatal(err)
	}

	var output string

	switch cfg.outputFlag {
	case "table":
		output = nbpCurrency.CreatePrettyOutput(cfg.langFlag)
	case "json", "xml":
		output = nbpCurrency.CreateRawOutput()
	case "csv":
		output = nbpCurrency.CreateCSVOutput(cfg.langFlag)
	}

	if cfg.clipFlag {
		err = clipboard.WriteAll(output)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println(output)
	}
}

// tableCommand - function for 'table' command (tables with exchange rates)
func tableCommand() {
	var err error

	err = checkArg("table", cfg.tableFlag, cfg.dateFlag, cfg.lastFlag, cfg.outputFlag, cfg.codeFlag)
	if err != nil {
		log.Fatal(err)
	}

	nbpTable := nbpapi.NewTable(cfg.tableFlag)
	if cfg.outputFlag == "xml" || cfg.outputFlag == "json" {
		err = nbpTable.TableRaw(cfg.dateFlag, cfg.lastFlag, cfg.repFormat)
	} else if cfg.lastFlag > 0 {
		err = nbpTable.TableLast(cfg.lastFlag)
	} else {
		err = nbpTable.TableByDate(cfg.dateFlag)
	}
	if err != nil {
		log.Fatal(err)
	}

	var output string

	switch cfg.outputFlag {
	case "table":
		output = nbpTable.CreatePrettyOutput(cfg.langFlag)
	case "json", "xml":
		output = nbpTable.CreateRawOutput()
	case "csv":
		output = nbpTable.CreateCSVOutput(cfg.langFlag)
	}

	if cfg.clipFlag {
		err = clipboard.WriteAll(output)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println(output)
	}
}
