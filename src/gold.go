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
