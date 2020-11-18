package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"
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

// getTableToday - funkcja zwraca json z tabelą kursów podanego typu na dziś (lub błąd)
func getTableToday(tableType string) ([]byte, error) {
	address := fmt.Sprintf(baseAddress+"/tables/%s/today/?format=json", tableType)
	return getJSON(address)
}

// getTableCurrent - funkcja zwraca bieżącą tabelę kursów walut danego typu
// (ostatnio opublikowaną tabelę danego typu) w formie json
func getTableCurrent(tableType string) ([]byte, error) {
	address := fmt.Sprintf(baseAddress+"/tables/%s/?format=json", tableType)
	return getJSON(address)
}

// getTableDay - funkcja zwraca tabelę kursów (json) danego typu dla podanego dnia (lub błąd)
func getTableDay(tableType string, day string) ([]byte, error) {
	address := fmt.Sprintf(baseAddress+"/tables/%s/%s/?format=json", tableType, day)
	return getJSON(address)
}

// getTableRange - funkcja zwraca tabele kursów danego typu dla podanego zakresu
// dat w formie json (lub błąd)
func getTableRange(tableType string, day string) ([]byte, error) {
	var startDate string
	var stopDate string

	temp := strings.Split(day, ":")
	if len(temp) != 2 {
		log.Fatal(errors.New("Nieprawidłowy format zakresu dat"))
	}

	startDate = temp[0]
	stopDate = temp[1]

	address := fmt.Sprintf(baseAddress+"/tables/%s/%s/%s/?format=json", tableType, startDate, stopDate)
	return getJSON(address)
}

// getTableLast - funkcja zwraca ostatnich n tabel kursów danego typu
// w formie json (lub błąd)
func getTableLast(tableType string, last string) ([]byte, error) {
	address := fmt.Sprintf(baseAddress+"/tables/%s/last/%s/?format=json", tableType, last)
	return getJSON(address)
}

// printTable - funkcja drukuje tabele kursów w konsoli
func printTable(result []byte, tableType string) {
	var nbpTables []exchangeTable
	var nbpTablesC []exchangeTableC

	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.Debug)

	if tableType != "C" {
		err := json.Unmarshal(result, &nbpTables)
		if err != nil {
			log.Fatal(err)
		}

		for _, item := range nbpTables {
			fmt.Println()
			fmt.Println("Typ tabeli:\t\t", item.Table)
			fmt.Println("Numer tabeli:\t\t", item.No)
			fmt.Println("Data publikacji:\t", item.EffectiveDate)
			fmt.Println()

			fmt.Fprintln(w, "KOD \t NAZWA \t ŚREDNI")
			fmt.Fprintln(w, "--- \t ----- \t -------")
			for _, currencyItem := range item.Rates {
				currencyValue := fmt.Sprintf("%.4f", currencyItem.Mid)
				fmt.Fprintln(w, currencyItem.Code+" \t "+currencyItem.Currency+" \t "+currencyValue)
			}

			w.Flush()
		}
	} else {
		err := json.Unmarshal(result, &nbpTablesC)
		if err != nil {
			log.Fatal(err)
		}

		for _, item := range nbpTablesC {
			fmt.Println()
			fmt.Println("Typ tabeli:\t\t", item.Table)
			fmt.Println("Numer tabeli:\t\t", item.No)
			fmt.Println("Data notowania:\t\t", item.TradingDate)
			fmt.Println("Data publikacji:\t", item.EffectiveDate)
			fmt.Println()

			fmt.Fprintln(w, "KOD \t NAZWA \t KUPNO \t SPRZEDAŻ ")
			fmt.Fprintln(w, "--- \t ----- \t ----- \t -------- ")
			for _, currencyItem := range item.Rates {
				currencyValueBid := fmt.Sprintf("%.4f", currencyItem.Bid)
				currencyValueAsk := fmt.Sprintf("%.4f", currencyItem.Ask)
				fmt.Fprintln(w, currencyItem.Code+" \t "+currencyItem.Currency+" \t "+currencyValueBid+" \t "+currencyValueAsk)
			}

			w.Flush()
		}
	}

	fmt.Println()
}

// printTableCSV - funkcja drukuje dane tabel kursów w konsoli w formie CSV
func printTableCSV(result []byte, tableType string) {
	var nbpTables []exchangeTable
	var nbpTablesC []exchangeTableC
	var tableNo string

	if tableType != "C" {
		err := json.Unmarshal(result, &nbpTables)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("TABELA,KOD,NAZWA,ŚREDNI")

		for _, item := range nbpTables {
			tableNo = item.No
			for _, currencyItem := range item.Rates {
				currencyValue := fmt.Sprintf("%.4f", currencyItem.Mid)
				fmt.Println(tableNo + "," + currencyItem.Code + "," + currencyItem.Currency + "," + currencyValue)
			}
		}
	} else {
		err := json.Unmarshal(result, &nbpTablesC)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("TABELA,KOD,NAZWA,KUPNO,SPRZEDAŻ")

		for _, item := range nbpTablesC {
			tableNo = item.No
			for _, currencyItem := range item.Rates {
				currencyValueBid := fmt.Sprintf("%.4f", currencyItem.Bid)
				currencyValueAsk := fmt.Sprintf("%.4f", currencyItem.Ask)
				fmt.Println(tableNo + "," + currencyItem.Code + "," + currencyItem.Currency + "," + currencyValueBid + "," + currencyValueAsk)
			}
			fmt.Println()
		}
	}

	fmt.Println()
}
