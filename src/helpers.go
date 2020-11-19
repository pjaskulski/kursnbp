package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

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

// inSlice - funkcja sprawdza czy podany string występuje we wskazanym wycinku
func inSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// checkArg - funkcja weryfikuje poprawność parametrów wywołania programu
func checkArg(cmd string, tFlag string, dFlag string, lFlag int, oFlag string, cFlag string) error {
	tableValues := []string{"A", "B", "C"}
	outputValues := []string{"table", "json", "csv"}
	currencyValuesA := []string{"THB", "USD", "AUD", "HKD", "CAD", "NZD", "SGD", "EUR", "HUF", "CHF",
		"GBP", "UAH", "JPY", "CZK", "DKK", "ISK", "NOK", "SEK", "HRK", "RON",
		"BGN", "TRY", "ILS", "CLP", "PHP", "MXN", "ZAR", "BRL", "MYR", "RUB",
		"IDR", "INR", "KRW", "CNY", "XDR"}
	currencyValuesB := []string{"MGA", "PAB", "ETB", "AFN", "VES", "BOB", "CRC", "SVC", "NIO", "GMD",
		"MKD", "DZD", "BHD", "IQD", "JOD", "KWD", "LYD", "RSD", "TND", "MAD",
		"AED", "STN", "BSD", "BBD", "BZD", "BND", "FJD", "GYD", "JMD", "LRD",
		"NAD", "SRD", "TTD", "XCD", "SBD", "ZWL", "VND", "AMD", "CVE", "AWG",
		"BIF", "XOF", "XAF", "XPF", "DJF", "GNF", "KMF", "CDF", "RWF", "EGP",
		"GIP", "LBP", "SSP", "SDG", "SYP", "GHS", "HTG", "PYG", "ANG", "PGK",
		"LAK", "MWK", "ZMW", "AOA", "MMK", "GEL", "MDL", "ALL", "HNL", "SLL",
		"SZL", "LSL", "AZN", "MZN", "NGN", "ERN", "TWD", "TMT", "MRU", "TOP",
		"MOP", "ARS", "DOP", "COP", "CUP", "UYU", "BWP", "GTQ", "IRR", "YER",
		"QAR", "OMR", "SAR", "KHR", "BYN", "LKR", "MVR", "MUR", "NPR", "PKR",
		"SCR", "PEN", "KGS", "TJS", "UZS", "KES", "SOS", "TZS", "UGX", "BDT",
		"WST", "KZT", "MNT", "VUV", "BAM"}
	currencyValuesC := []string{"USD", "AUD", "CAD", "EUR", "HUF", "GBP", "JPY", "CZK", "DKK", "NOK",
		"SEK", "XDR"}

	// output
	if oFlag == "" {
		return errors.New("Brak wartości parametru --output, należy podać format danych wyjściowych")
	}
	if !inSlice(outputValues, oFlag) {
		return errors.New("Nieprawidłowa wartość parametru --output, dozwolone: table, json, csv")
	}

	// last
	if lFlag == 0 && dFlag == "" {
		return errors.New("Należy podać wartość jednego z parametrów: --date lub --last")
	}
	if lFlag < 0 {
		return errors.New("Nieprawidłowa wartość parametru --last, dozwolona wartość > 0")
	}
	if lFlag > 0 && dFlag != "" {
		return errors.New("Należy podać wartość tylko jednego z parametrów: albo --date albo --last")
	}

	// date
	if lFlag == 0 && dFlag != "" {
		var isValid bool = true

		if dFlag != "today" && dFlag != "current" {
			if len(dFlag) != 10 && len(dFlag) != 21 {
				isValid = false
			}
			if len(dFlag) == 10 {
				re10 := regexp.MustCompile("\\d{4}-\\d{2}-\\d{2}")
				if !re10.MatchString(dFlag) {
					isValid = false
				}
			}
			if len(dFlag) == 21 {
				re21 := regexp.MustCompile("\\d{4}-\\d{2}-\\d{2}\\:\\d{4}-\\d{2}-\\d{2}")
				if !re21.MatchString(dFlag) {
					isValid = false
				}
			}
			if !isValid {
				return errors.New(`Nieprawidłowa wartość parametru --date, dozwolone wartości: 
'today', 'current', 'RRRR-MM-DD' lub 'RRRR-MM-DD:RRRR-MM-DD'`)
			}
		}
	}

	// table
	if cmd == "table" {
		if tFlag == "" {
			return errors.New("Brak wartości parametru --table, należy podać typ tabeli kursów")
		}
		if !inSlice(tableValues, tFlag) {
			return errors.New("Nieprawidłowa wartość parametru --table, dozwolone: A, B lub C")
		}
	}

	// currency
	var errMessage string

	if cmd == "currency" {
		if cFlag == "" {
			return errors.New("Brak wartości parametru --code, należy podać kod waluty")
		}
		if tFlag == "" {
			return errors.New("Brak wartości parametru --table, należy podać typ tabeli kursów")
		}

		if !inSlice(tableValues, tFlag) {
			return errors.New("Nieprawidłowa wartość parametru --table, dozwolone: A, B lub C")
		}

		if tFlag == "A" {
			if !inSlice(currencyValuesA, cFlag) {
				errMessage = "Nieprawidłowa wartość parametru --code, "
				errMessage += "dozwolony poprawny kod waluty z dostępnych dla tabeli A: "
				errMessage += strings.Join(currencyValuesA, ", ")
				return errors.New(errMessage)
			}
		} else if tFlag == "B" {
			if !inSlice(currencyValuesB, cFlag) {
				errMessage = "Nieprawidłowa wartość parametru --code, "
				errMessage += "dozwolony poprawny kod waluty z dostępnych dla tabeli B: "
				errMessage += strings.Join(currencyValuesB, ", ")
				return errors.New(errMessage)
			}
		} else if tFlag == "C" {
			if !inSlice(currencyValuesC, cFlag) {
				errMessage = "Nieprawidłowa wartość parametru --code, "
				errMessage += "dozwolony poprawny kod waluty z dostępnych dla tabeli C: "
				errMessage += strings.Join(currencyValuesC, ", ")
				return errors.New(errMessage)
			}
		}
	}

	return nil
}
