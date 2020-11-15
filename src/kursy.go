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

// funkcja sprawdza czy w przekazanym stringu znajdują się tylko znaki dozwolone
// dla parametrów typu data lub zakres dat
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

// funkcja zwraca string zawierający tabelę kursów podanego typu na dziś (lub błąd)
func getTodayTable(tableType string) ([]byte, error) {
	address := fmt.Sprintf("http://api.nbp.pl/api/exchangerates/tables/%s/today/?format=json", tableType)
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

// funkcja zwraca bieżącą wartość tabeli kursów danego typu (ostatnio opublikowaną tabelę
// danego typu)
func getCurrentTable(tableType string) ([]byte, error) {
	address := fmt.Sprintf("http://api.nbp.pl/api/exchangerates/tables/%s/?format=json", tableType)
	r, err := http.Get(address)
	if err != nil {
		log.Fatal(err)
	}

	defer r.Body.Close()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	return data, nil
}

// funkcja zwraca tabelę kursów danego typu dla podanego dnia (lub błąd)
func getDayTable(tableType string, day string) ([]byte, error) {
	address := fmt.Sprintf("http://api.nbp.pl/api/exchangerates/tables/%s/%s/?format=json", tableType, day)
	r, err := http.Get(address)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	if r.StatusCode == 404 {
		errorText := fmt.Sprintf("Nie znaleziono tabeli kursów opublikowanej w żądanym dniu: %s", day)
		return nil, errors.New(errorText)
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	return data, nil
}

// funkcja zwraca tabele kursów danego typu dla podanego zakresu dat (lub błąd)
func getRangeTable(tableType string, day string) ([]byte, error) {
	var startDate string
	var stopDate string

	temp := strings.Split(day, ":")
	if len(temp) != 2 {
		log.Fatal(errors.New("Nieprawidłowy format zakresu dat"))
	}

	startDate = temp[0]
	stopDate = temp[1]

	address := fmt.Sprintf("http://api.nbp.pl/api/exchangerates/tables/%s/%s/%s/?format=json", tableType, startDate, stopDate)
	r, err := http.Get(address)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	if r.StatusCode == 404 {
		errorText := fmt.Sprintf("Nie znaleziono tabeli kursów dla żądanego zakresu: %s", day)
		return nil, errors.New(errorText)
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	return data, nil
}

// funkcja zwraca ostatnie tabele kursów danego typu dla  (lub błąd)
func getLastTable(tableType string, last string) ([]byte, error) {
	address := fmt.Sprintf("http://api.nbp.pl/api/exchangerates/tables/%s/last/%s/?format=json", tableType, last)
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

// funkcja drukuje tabele kursów w konsoli
func printTable(result []byte) {
	var nbpTables []currencyTable
	err := json.Unmarshal(result, &nbpTables)
	if err != nil {
		log.Fatal(err)
	}
	// druk tabeli z kursami w oknie konsoli
	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.Debug)

	for _, item := range nbpTables {
		fmt.Println()
		fmt.Println("Typ tabeli:", item.Table)
		fmt.Println("Numer tabeli:", item.No)
		fmt.Println("Data publikacji:", item.EffectiveDate)
		fmt.Println()

		for _, currencyItem := range item.Rates {
			currencyValue := fmt.Sprintf("%.4f", currencyItem.Mid)
			fmt.Fprintln(w, currencyItem.Code+" \t "+currencyItem.Currency+" \t "+currencyValue)
		}

		w.Flush()
	}

	fmt.Println()
}
