package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type rateTable struct {
	Currency string  `json:"currency"`
	Code     string  `json:"code"`
	Mid      float64 `json:"mid"`
}

type exchangeTable struct {
	Table         string      `json:"table"`
	No            string      `json:"no"`
	EffectiveDate string      `json:"effectiveDate"`
	Rates         []rateTable `json:"rates"`
}

type rateTableC struct {
	Currency string  `json:"currency"`
	Code     string  `json:"code"`
	Bid      float64 `json:"bid"`
	Ask      float64 `json:"ask"`
}

type exchangeTableC struct {
	Table         string       `json:"table"`
	No            string       `json:"no"`
	TradingDate   string       `json:"tradingDate"`
	EffectiveDate string       `json:"effectiveDate"`
	Rates         []rateTableC `json:"rates"`
}

type rateCurrency struct {
	No            string  `json:"no"`
	EffectiveDate string  `json:"effectiveDate"`
	Mid           float64 `json:"mid"`
}

type exchangeCurrency struct {
	Table    string         `json:"table"`
	Currency string         `json:"currency"`
	Code     string         `json:"code"`
	Rates    []rateCurrency `json:"rates"`
}

type rateCurrencyC struct {
	No            string  `json:"no"`
	EffectiveDate string  `json:"effectiveDate"`
	Bid           float64 `json:"bid"`
	Ask           float64 `json:"ask"`
}

type exchangeCurrencyC struct {
	Table    string          `json:"table"`
	Currency string          `json:"currency"`
	Code     string          `json:"code"`
	Rates    []rateCurrencyC `json:"rates"`
}

type rateGold struct {
	Data string  `json:"data"`
	Cena float64 `json:"cena"`
}

// argumenty startowe:
// -table <type> - typ tabeli kursów (A, B, lub C)
// -day=today - kurs na dziś
// -day=<date> - kurs na dzień (format RRRR-MM-DD)
// -day=<startDate>:<endDate> - kursy z zakresu dat (format RRRR-MM-DD:RRRR-MM-DD)
// -last=<number> - ostatnich <number> kursów
// -out=<output> - format wyjścia, domyślnie json (json, table - tabela tekstowa)
// np.
// kursnbp -table A -day today
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

	fmt.Println("Kursy NBP - klient tekstowy")

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

	if outputFlag == "json" {
		fmt.Println(string(result))
	} else if outputFlag == "table" {
		if currencyFlag == "ALL" {
			if tableFlag != "C" {
				printTable(result)
			} else {
				printTableC(result)
			}
		} else if currencyFlag == "GOLD" {
			printGold(result)
		} else {
			if tableFlag != "C" {
				printCurrency(result)
			} else {
				printCurrencyC(result)
			}
		}
	}
}
