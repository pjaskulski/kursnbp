package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type rate struct {
	Currency string  `json:"currency"`
	Code     string  `json:"code"`
	Mid      float64 `json:"mid"`
}

type currencyTable struct {
	Table         string `json:"table"`
	No            string `json:"no"`
	EffectiveDate string `json:"effectiveDate"`
	Rates         []rate `json:"rates"`
}

// argumenty startowe:
// 	-table <type> - typ tabeli kursów (A, B, lub C)
// 	-today - kurs na dziś
// 	-day <date> - kurs na dzień (format RRRR-MM-DD)
//  -day <startDate>:<endDate> - kursy z zakresu dat (format RRRR-MM-DD:RRRR-MM-DD)
// 	-last <number> - ostatnich <number> kursów
//  -out <output> - format wyjścia, domyślnie json (json, table - tabela tekstowa)
//  np.
//  kursnbp -table A -day today
func main() {
	var tableFlag string
	var dayFlag string
	var outputFlag string
	var lastFlag string

	flag.StringVar(&tableFlag, "table", "", "typ tabeli kursów (A, B lub C)")
	flag.StringVar(&dayFlag, "day", "", "data tabeli kursów (RRRR-MM-DD lub: today, current lub zakres dat RRRR-MM-DD:RRRR-MM-DD)")
	flag.StringVar(&outputFlag, "output", "json", "format wyjścia (json, table)")
	flag.StringVar(&lastFlag, "last", "", "seria ostatnich tabel kursów")
	flag.Parse()

	fmt.Println("Kursy NBP - klient tekstowy")

	if tableFlag == "" || (dayFlag == "" && lastFlag == "") {
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
		result, err = getLastTable(tableFlag, lastFlag)
		if err != nil {
			log.Fatal(err)
		}

	} else {
		if dayFlag == "today" {
			result, err = getTodayTable(tableFlag)
			if err != nil {
				log.Fatal(err)
			}

		} else if dayFlag == "current" {
			result, err = getCurrentTable(tableFlag)
			if err != nil {
				log.Fatal(err)
			}

		} else if len(dayFlag) == 10 && charAllowed(dayFlag, false) {
			result, err = getDayTable(tableFlag, dayFlag)
			if err != nil {
				log.Fatal(err)
			}

		} else if len(dayFlag) == 21 && charAllowed(dayFlag, true) {
			result, err = getRangeTable(tableFlag, dayFlag)
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
		printTable(result)
	}
}
