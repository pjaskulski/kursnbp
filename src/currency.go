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

// getCurrencyLast - funkcja zwraca ostatnich n kursów waluty danego typu
// w formie json (lub błąd)
func getCurrencyLast(tableType string, last string, currency string) ([]byte, error) {
	address := fmt.Sprintf(baseAddress+"/rates/%s/%s/last/%s/?format=json", tableType, currency, last)
	return getJSON(address)
}

// getCurrencyToday - funkcja zwraca json z kursem waluty podanego typu na dziś (lub błąd)
func getCurrencyToday(tableType string, currency string) ([]byte, error) {
	address := fmt.Sprintf(baseAddress+"/rates/%s/%s/today/?format=json", tableType, currency)
	return getJSON(address)
}

// getCurrencyCurrent - funkcja zwraca bieżący kurs waluty danego typu (ostatnio
// opublikowany kurs waluty danego typu) w formie json
func getCurrencyCurrent(tableType string, currency string) ([]byte, error) {
	address := fmt.Sprintf(baseAddress+"/rates/%s/%s/?format=json", tableType, currency)
	return getJSON(address)
}

// getCurrencyDay - funkcja zwraca kurs waluty (json) danego typu dla podanego dnia (lub błąd)
func getCurrencyDay(tableType string, day string, currency string) ([]byte, error) {
	address := fmt.Sprintf(baseAddress+"/rates/%s/%s/%s/?format=json", tableType, currency, day)
	return getJSON(address)
}

// getCurrencyRange - funkcja zwraca tabele kursów danego typu dla podanego
// zakresu dat w formie json (lub błąd)
func getCurrencyRange(tableType string, day string, currency string) ([]byte, error) {
	var startDate string
	var stopDate string

	temp := strings.Split(day, ":")
	if len(temp) != 2 {
		log.Fatal(errors.New("Nieprawidłowy format zakresu dat"))
	}

	startDate = temp[0]
	stopDate = temp[1]

	address := fmt.Sprintf(baseAddress+"/rates/%s/%s/%s/%s/?format=json", tableType, currency, startDate, stopDate)
	return getJSON(address)
}

// printCurrency - funkcja drukuje kursy waluty w konsoli
func printCurrency(result []byte, tableType string) {
	var nbpCurrency exchangeCurrency
	var nbpCurrencyC exchangeCurrencyC

	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.Debug)

	if tableType != "C" {
		err := json.Unmarshal(result, &nbpCurrency)
		if err != nil {
			log.Fatal(err)
		}

		// druk kursów waluty w oknie konsoli
		fmt.Println()
		fmt.Println("Typ tabeli:\t", nbpCurrency.Table)
		fmt.Println("Nazwa waluty:\t", nbpCurrency.Currency)
		fmt.Println("Kod waluty:\t", nbpCurrency.Code)
		fmt.Println()

		fmt.Fprintln(w, "TABELA \t DATA \t ŚREDNI")
		fmt.Fprintln(w, "------ \t ---- \t -------")
		for _, currencyItem := range nbpCurrency.Rates {
			currencyValue := fmt.Sprintf("%.4f", currencyItem.Mid)
			fmt.Fprintln(w, currencyItem.No+" \t "+currencyItem.EffectiveDate+" \t "+currencyValue)
		}
	} else {
		err := json.Unmarshal(result, &nbpCurrencyC)
		if err != nil {
			log.Fatal(err)
		}
		// druk kursów waluty w oknie konsoli
		fmt.Println()
		fmt.Println("Typ tabeli:\t", nbpCurrencyC.Table)
		fmt.Println("Nazwa waluty:\t", nbpCurrencyC.Currency)
		fmt.Println("Kod waluty:\t", nbpCurrencyC.Code)
		fmt.Println()

		fmt.Fprintln(w, "TABELA \t DATA \t KUPNO \t SPRZEDAŻ ")
		fmt.Fprintln(w, "------ \t ---- \t ----- \t -------- ")
		for _, currencyItem := range nbpCurrencyC.Rates {
			currencyValueBid := fmt.Sprintf("%.4f", currencyItem.Bid)
			currencyValueAsk := fmt.Sprintf("%.4f", currencyItem.Ask)
			fmt.Fprintln(w, currencyItem.No+" \t "+currencyItem.EffectiveDate+" \t "+currencyValueBid+" \t "+currencyValueAsk)
		}
	}
	w.Flush()

	fmt.Println()
}

// printCurrency - funkcja drukuje kursy waluty w konsoli
func printCurrencyCSV(result []byte, tableType string) {
	var nbpCurrency exchangeCurrency
	var nbpCurrencyC exchangeCurrencyC

	if tableType != "C" {
		err := json.Unmarshal(result, &nbpCurrency)
		if err != nil {
			log.Fatal(err)
		}

		// druk kursów waluty w oknie konsoli
		fmt.Println("TABELA,DATA,ŚREDNI")
		for _, currencyItem := range nbpCurrency.Rates {
			currencyValue := fmt.Sprintf("%.4f", currencyItem.Mid)
			fmt.Println(currencyItem.No + "," + currencyItem.EffectiveDate + "," + currencyValue)
		}
	} else {
		err := json.Unmarshal(result, &nbpCurrencyC)
		if err != nil {
			log.Fatal(err)
		}

		// druk danych kursów waluty w oknie konsoli w formie CSV
		fmt.Println("TABELA,DATA,KUPNO,SPRZEDAŻ")
		for _, currencyItem := range nbpCurrencyC.Rates {
			currencyValueBid := fmt.Sprintf("%.4f", currencyItem.Bid)
			currencyValueAsk := fmt.Sprintf("%.4f", currencyItem.Ask)
			fmt.Println(currencyItem.No + "," + currencyItem.EffectiveDate + "," + currencyValueBid + "," + currencyValueAsk)
		}
	}
	fmt.Println()
}
