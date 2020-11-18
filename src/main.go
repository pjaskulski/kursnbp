package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// adresy bazowe dla kursów walut i cen złota
const (
	baseAddress     string = "http://api.nbp.pl/api/exchangerates"
	baseAddressGold string = "http://api.nbp.pl/api/cenyzlota"
)

var version string = "0.1.0"

// kursnbp - command line tool for downloading exchange rates and gold prices
// from the website of the National Bank of Poland
func main() {
	var tableFlag string
	var dayFlag string
	var outputFlag string
	var lastFlag string
	var currencyFlag string

	flag.StringVar(&tableFlag, "table", "", "typ tabeli kursów (A, B lub C), dla cen złota ignorowany")
	flag.StringVar(&dayFlag, "day", "", "data tabeli kursów (RRRR-MM-DD lub: today, current lub zakres dat RRRR-MM-DD:RRRR-MM-DD)")
	flag.StringVar(&outputFlag, "output", "table", "format wyjścia (json, table)")
	flag.StringVar(&lastFlag, "last", "", "seria ostatnich tabel kursów, ostatnich kursów waluty lub złota")
	flag.StringVar(&currencyFlag, "currency", "ALL", "kod waluty lub ALL = cała tabela kursów lub GOLD - cena złota")
	flag.Parse()

	if outputFlag == "table" {
		fmt.Println("Kursy NBP - klient tekstowy, wersja " + version)
	}

	if (currencyFlag != "GOLD" && tableFlag == "") || (dayFlag == "" && lastFlag == "") {
		fmt.Println("Parametry wywołania programu:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if dayFlag != "" && lastFlag != "" {
		fmt.Println("Parametry -last i -day nie mogą występować jednocześnie")
		os.Exit(1)
	}

	if !strings.Contains("ABC", tableFlag) {
		log.Fatal("Należy podać poprawny typ tabeli kursów: A, B lub C")
	}

	var result []byte
	var err error

	if lastFlag != "" {
		if currencyFlag == "ALL" {
			result, err = getTableLast(tableFlag, lastFlag)
		} else if currencyFlag == "GOLD" {
			result, err = getGoldLast(lastFlag)
		} else {
			result, err = getCurrencyLast(tableFlag, lastFlag, currencyFlag)
		}
		if err != nil {
			log.Fatal(err)
		}

	} else {
		if dayFlag == "today" {
			if currencyFlag == "ALL" {
				result, err = getTableToday(tableFlag)
			} else if currencyFlag == "GOLD" {
				result, err = getGoldToday()
			} else {
				result, err = getCurrencyToday(tableFlag, currencyFlag)
			}
			if err != nil {
				log.Fatal(err)
			}

		} else if dayFlag == "current" {
			if currencyFlag == "ALL" {
				result, err = getTableCurrent(tableFlag)
			} else if currencyFlag == "GOLD" {
				result, err = getGoldCurrent()
			} else {
				result, err = getCurrencyCurrent(tableFlag, currencyFlag)
			}
			if err != nil {
				log.Fatal(err)
			}

		} else if len(dayFlag) == 10 && charAllowed(dayFlag, false) {
			if currencyFlag == "ALL" {
				result, err = getTableDay(tableFlag, dayFlag)
			} else if currencyFlag == "GOLD" {
				result, err = getGoldDay(dayFlag)
			} else {
				result, err = getCurrencyDay(tableFlag, dayFlag, currencyFlag)
			}
			if err != nil {
				log.Fatal(err)
			}

		} else if len(dayFlag) == 21 && charAllowed(dayFlag, true) {
			if currencyFlag == "ALL" {
				result, err = getTableRange(tableFlag, dayFlag)
			} else if currencyFlag == "GOLD" {
				result, err = getGoldRange(dayFlag)
			} else {
				result, err = getCurrencyRange(tableFlag, dayFlag, currencyFlag)
			}
			if err != nil {
				log.Fatal(err)
			}

		} else {
			log.Fatal("Nieprawidłowy format parametru -day")
		}
	}

	// output
	if outputFlag == "json" {
		fmt.Println(string(result))
	} else if outputFlag == "table" {
		if currencyFlag == "ALL" {
			printTable(result, tableFlag)
		} else if currencyFlag == "GOLD" {
			printGold(result)
		} else {
			printCurrency(result, tableFlag)
		}
	} else if outputFlag == "csv" {
		if currencyFlag == "ALL" {
			printTableCSV(result, tableFlag)
		} else if currencyFlag == "GOLD" {
			printGoldCSV(result)
		} else {
			printCurrencyCSV(result, tableFlag)
		}
	}
}
