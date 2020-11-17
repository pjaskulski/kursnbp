package main

import (
	"encoding/json"
	"log"
	"strings"
	"testing"
	"time"
)

// aby nie męczyć serwera NBP za bardzo...
func littleDelay() {
	time.Sleep(time.Millisecond * 500)
}

func TestCharactersAllowed(t *testing.T) {

	tests := []struct {
		name      string
		text      string
		dateRange bool
		want      bool
	}{
		{
			name:      "Poprawna data",
			text:      "2020-12-11",
			dateRange: false,
			want:      true,
		},
		{
			name:      "Niepoprawna format daty",
			text:      "2020/11/10",
			dateRange: false,
			want:      false,
		},
		{
			name:      "Niepoprawny format daty 2",
			text:      "2020:11:10",
			dateRange: false,
			want:      false,
		},
		{
			name:      "Poprawny zakres",
			text:      "2020-11-12:2020-11-13",
			dateRange: true,
			want:      true,
		},
		{
			name:      "Niepoprawny zakres",
			text:      "2020/11/12-2020/11/13",
			dateRange: true,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := charAllowed(tt.text, tt.dateRange)
			if result != tt.want {
				t.Errorf("oczekiwano: %t; otrzymano: %t", tt.want, result)
			}
		})
	}

}

func TestGetCurrencyCurrent(t *testing.T) {
	var table string = "A"
	var currency string = "CHF"

	littleDelay()
	result, err := getCurrencyCurrent(table, currency)
	if err != nil {
		t.Errorf("oczekiwano err == nil, otrzymano err != nil")
	}
	if !strings.Contains(string(result), "\"table\":\"A\",\"currency\":\"frank szwajcarski\"") {
		t.Errorf("otrzymano niepoprawną zawartość json")
	}
}

func TestGetCurrencyCurrentXXX(t *testing.T) {
	var table string = "A"
	var currency string = "XXX" // niepoprawny kod waluty

	littleDelay()
	_, err := getCurrencyCurrent(table, currency)
	if err == nil {
		t.Errorf("oczekiwano err != nil, otrzymano err == nil")
	}
}

func TestGetCurrencyDay(t *testing.T) {
	var table string = "A"
	var currency string = "CHF"
	var day string = "2020-11-13" // Friday - ok, kurs CHF = 4.1605

	littleDelay()
	result, err := getCurrencyDay(table, day, currency)
	if err != nil {
		t.Errorf("oczekiwano err == nil, otrzymano err != nil")
	}
	// mało eleganckie ale skuteczne
	if !strings.Contains(string(result), "\"effectiveDate\":\"2020-11-13\",\"mid\":4.1605") {
		t.Errorf("niepoprawna zawartość json, kurs CHF 13.11.2020 wynosił 4.1605")
	}
}

func TestGetCurrencyDaySaturday(t *testing.T) {
	var table string = "A"
	var currency string = "CHF"
	var day string = "2020-11-14" // Saturday - no table of exchange rates

	littleDelay()
	_, err := getCurrencyDay(table, day, currency)
	if err == nil {
		t.Errorf("oczekiwano err != nil, otrzymano err == nil")
	}
}

func TestGetCurrencyLast(t *testing.T) {
	var table string = "A"
	var currency string = "CHF"
	var last string = "5"

	littleDelay()
	_, err := getCurrencyLast(table, last, currency)
	if err != nil {
		t.Errorf("oczekiwano err == nil, otrzymano err != nil")
	}
}

func TestGetCurrencyLastFailed(t *testing.T) {
	var table string = "A"
	var currency string = "CHF"
	var last string = "500" // za dużo kursów, max = 255

	littleDelay()
	_, err := getCurrencyLast(table, last, currency)
	if err == nil {
		t.Errorf("oczekiwano err != nil, otrzymano err == nil")
	}
}

func TestGetCurrencyRange(t *testing.T) {
	var table string = "A"
	var currency string = "CHF"
	var day string = "2020-11-12:2020-11-13" // poprawny zakres dat, spodziewane 2 kursy

	littleDelay()
	result, err := getCurrencyRange(table, day, currency)
	if err != nil {
		t.Errorf("oczekiwano err == nil, otrzymano err != nil")
	}

	var nbpCurrency exchangeCurrency
	err = json.Unmarshal(result, &nbpCurrency)
	if err != nil {
		log.Fatal(err)
	}
	var ratesCount int = len(nbpCurrency.Rates)
	if ratesCount != 2 {
		t.Errorf("oczekiwana liczba kursów == 2, otrzymano %d", ratesCount)
	}
}

func TestGetCurrencyRangeFailed(t *testing.T) {
	var table string = "A"
	var currency string = "CHF"
	var day string = "2020-11-12:2020-11-10" // niepoprawny zakres dat

	littleDelay()
	_, err := getCurrencyRange(table, day, currency)
	if err == nil {
		t.Errorf("oczekiwano err != nil, otrzymano err == nil")
	}
}

func TestGetCurrencyToday(t *testing.T) {
	var table string = "A"
	var currency string = "CHF"
	today := time.Now()
	var day string = today.Format("2006-01-02")

	littleDelay()
	_, err := getCurrencyDay(table, day, currency)
	if err == nil {
		_, err := getCurrencyToday(table, currency)
		if err != nil {
			t.Errorf("oczekiwano err == nil, otrzymano err != nil")
		}
	}
}

func TestGetTableCurrent(t *testing.T) {
	var table string = "A"

	littleDelay()
	result, err := getTableCurrent(table)
	if err != nil {
		t.Errorf("oczekiwano err == nil, otrzymano err != nil")
	}
	if !strings.Contains(string(result), "{\"table\":\"A\",\"no\":") {
		t.Errorf("otrzymano niepoprawną zawartość json")
	}
}

func TestGetTableDay(t *testing.T) {
	var table string = "A"
	var day string = "2020-11-17"
	var tableNo string = "224/A/NBP/2020"

	littleDelay()
	result, err := getTableDay(table, day)
	if err != nil {
		t.Errorf("oczekiwano err == nil, otrzymano err != nil")
	}
	if !json.Valid(result) {
		t.Errorf("otrzymano niepoprawną zawartość json")
	}

	var nbpTables []exchangeTable
	err = json.Unmarshal(result, &nbpTables)
	if err != nil {
		log.Fatal(err)
	}

	if nbpTables[0].Table != table {
		t.Errorf("niepoprawny typ tabeli, oczekiwano %s, otrzymano %s", table, nbpTables[0].Table)
	}
	if nbpTables[0].No != tableNo {
		t.Errorf("niepoprawny numer tabeli, oczekiwano %s, otrzymano %s", tableNo, nbpTables[0].No)
	}
	if nbpTables[0].EffectiveDate != day {
		t.Errorf("niepoprawna data publikacji, oczekiwano %s, otrzymano %s", day, nbpTables[0].EffectiveDate)
	}
}

func TestGetTableRange(t *testing.T) {
	var table string = "A"
	var day string = "2020-11-16:2020-11-17"

	littleDelay()
	result, err := getTableRange(table, day)
	if err != nil {
		t.Errorf("oczekiwano err == nil, otrzymano err != nil")
	}
	if !json.Valid(result) {
		t.Errorf("otrzymano niepoprawną zawartość json")
	}

	var nbpTables []exchangeTable
	err = json.Unmarshal(result, &nbpTables)
	if err != nil {
		log.Fatal(err)
	}

	if len(nbpTables) != 2 {
		t.Errorf("oczekiwano 2 tabel kursów, otrzymano %d", len(nbpTables))
	}

	if nbpTables[0].Table != table {
		t.Errorf("niepoprawny typ tabeli, oczekiwano %s, otrzymano %s", table, nbpTables[0].Table)
	}

	if nbpTables[1].Table != table {
		t.Errorf("niepoprawny typ tabeli, oczekiwano %s, otrzymano %s", table, nbpTables[1].Table)
	}
}

func TestGetTableLast(t *testing.T) {
	var table string = "A"
	var lastNo string = "5"

	littleDelay()
	result, err := getTableLast(table, lastNo)
	if err != nil {
		t.Errorf("oczekiwano err == nil, otrzymano err != nil")
	}
	if !json.Valid(result) {
		t.Errorf("otrzymano niepoprawną zawartość json")
	}

	var nbpTables []exchangeTable
	err = json.Unmarshal(result, &nbpTables)
	if err != nil {
		log.Fatal(err)
	}

	if len(nbpTables) != 5 {
		t.Errorf("oczekiwano 5 tabel kursów, otrzymano %d", len(nbpTables))
	}
}

func TestGetTableToday(t *testing.T) {
	var table string = "A"
	today := time.Now()
	var day string = today.Format("2006-01-02")

	littleDelay()
	_, err := getTableDay(table, day)
	if err != nil {
		_, err := getTableToday(table)
		if err != nil {
			t.Errorf("oczekiwano err == nil, otrzymano err != nil")
		}
	}
}

func TestGetTableTodayFailed(t *testing.T) {
	var table string = "D"

	littleDelay()
	_, err := getTableToday(table)
	if err == nil {
		t.Errorf("oczekiwano err != nil, otrzymano err != nil")
	}
}
