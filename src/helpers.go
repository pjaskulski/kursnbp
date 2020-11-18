package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

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

// aby nie męczyć serwera NBP za bardzo...
func littleDelay() {
	time.Sleep(time.Millisecond * 500)
}
