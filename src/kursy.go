package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"
)

// adresy bazowe dla kursów walut i cen złota
var baseAddress string = "http://api.nbp.pl/api/exchangerates"
var baseAddressGold string = "http://api.nbp.pl/api/cenyzlota"

// charAllowed - funkcja sprawdza czy w przekazanym stringu znajdują się tylko
// znaki dozwolone dla parametrów typu data lub zakres dat
func charAllowed(text string, dateRange bool) bool {
	var characters = "0123456789-"
	var result bool = true

	if dateRange {
		characters += ":"
	}

	for _, item := range text {
		if !strings.Contains(characters, string(item)) {
			result = false
			break
		}
	}
	return result
}

// getJSON - uniwersalna funkcja zwracająca json (lub błąd) na podstawie przekazanego adresu
func getJSON(address string) ([]byte, error) {
	r, err := http.Get(address)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	if r.StatusCode >= 400 {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		return nil, errors.New(string(body))
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	return data, nil
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
func printTable(result []byte) {
	var nbpTables []exchangeTable
	err := json.Unmarshal(result, &nbpTables)
	if err != nil {
		log.Fatal(err)
	}
	// druk tabeli z kursami w oknie konsoli
	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.Debug)

	for _, item := range nbpTables {
		fmt.Println()
		fmt.Println("Typ tabeli:\t\t", item.Table)
		fmt.Println("Numer tabeli:\t\t", item.No)
		fmt.Println("Data publikacji:\t", item.EffectiveDate)
		fmt.Println()

		fmt.Fprintln(w, "KOD \t NAZWA \t WARTOŚĆ")
		fmt.Fprintln(w, "--- \t ----- \t -------")
		for _, currencyItem := range item.Rates {
			currencyValue := fmt.Sprintf("%.4f", currencyItem.Mid)
			fmt.Fprintln(w, currencyItem.Code+" \t "+currencyItem.Currency+" \t "+currencyValue)
		}

		w.Flush()
	}

	fmt.Println()
}

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
func printCurrency(result []byte) {
	var nbpCurrency exchangeCurrency
	err := json.Unmarshal(result, &nbpCurrency)
	if err != nil {
		log.Fatal(err)
	}
	// druk kursów waluty w oknie konsoli
	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.Debug)

	fmt.Println()
	fmt.Println("Typ tabeli:\t", nbpCurrency.Table)
	fmt.Println("Nazwa waluty:\t", nbpCurrency.Currency)
	fmt.Println("Kod waluty:\t", nbpCurrency.Code)
	fmt.Println()

	fmt.Fprintln(w, "TABELA \t DATA \t WARTOŚĆ")
	fmt.Fprintln(w, "------ \t ---- \t -------")
	for _, currencyItem := range nbpCurrency.Rates {
		currencyValue := fmt.Sprintf("%.4f", currencyItem.Mid)
		fmt.Fprintln(w, currencyItem.No+" \t "+currencyItem.EffectiveDate+" \t "+currencyValue)
	}
	w.Flush()

	fmt.Println()
}

// getGoldToday - funkcja zwraca dzisiejszą cenę złota
// w formie json, lub błąd
func getGoldToday() ([]byte, error) {
	address := fmt.Sprintf(baseAddressGold + "/today?format=json")
	return getJSON(address)
}

// getGoldCurrent - funkcja zwraca bieżącą cenę złota
// (ostatnio opublikowaną cenę) w formie json, lub błąd
func getGoldCurrent() ([]byte, error) {
	address := fmt.Sprintf(baseAddressGold + "?format=json")
	return getJSON(address)
}

// getGoldLast - funkcja zwraca ostatnich n cen złota
// w formie json, lub błąd
func getGoldLast(last string) ([]byte, error) {
	address := fmt.Sprintf(baseAddressGold+"/last/%s?format=json", last)
	return getJSON(address)
}

// getGoldDay - funkcja zwraca cenę złota w podanym dniu
// w formie json, lub błąd
func getGoldDay(day string) ([]byte, error) {
	address := fmt.Sprintf(baseAddressGold+"/%s?format=json", day)
	return getJSON(address)
}

// getGoldRange - funkcja zwraca ceny złota w podanym zakresie dat
// w formie json, lub błąd
func getGoldRange(day string) ([]byte, error) {
	var startDate string
	var stopDate string

	temp := strings.Split(day, ":")
	if len(temp) != 2 {
		log.Fatal(errors.New("Nieprawidłowy format zakresu dat"))
	}

	startDate = temp[0]
	stopDate = temp[1]

	address := fmt.Sprintf(baseAddressGold+"/%s/%s?format=json", startDate, stopDate)
	return getJSON(address)
}

// printGold - funkcja drukuje ceny złota w konsoli
func printGold(result []byte) {
	var nbpGold []rateGold
	err := json.Unmarshal(result, &nbpGold)
	if err != nil {
		log.Fatal(err)
	}
	// druk cen złota w oknie konsoli
	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.Debug)

	fmt.Println()
	fmt.Println("Cena złota (1 g złota w próbie 1000)")
	fmt.Println()

	fmt.Fprintln(w, "DATA \t CENA")
	fmt.Fprintln(w, "---- \t ---- ")
	for _, goldItem := range nbpGold {
		goldValue := fmt.Sprintf("%.4f", goldItem.Cena)
		fmt.Fprintln(w, goldItem.Data+" \t "+goldValue)
	}
	w.Flush()

	fmt.Println()
}
