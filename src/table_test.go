package main

import (
	"encoding/json"
	"log"
	"testing"
	"time"
)

func TestGetTableCurrent(t *testing.T) {
	var table string = "A"

	littleDelay()
	result, err := getTableCurrent(table)
	if err != nil {
		t.Errorf("oczekiwano err == nil, otrzymano err != nil")
	}
	if !json.Valid(result) {
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
	if err == nil {
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
