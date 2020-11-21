// auxiliary program functions

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

// getJSON - universal function that returns JSON (or error) based on the address provided
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

// littleDelay - delay function, so as not to bother the NBP server too much...
func littleDelay() {
	time.Sleep(time.Millisecond * 500)
}

// inSlice - function checks if the specified string is present in the specified slice
func inSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// checkArg - function verifies the correctness of program call parameters
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
		return errors.New(l.Get("No --output parameter value, output format must be specified"))
	}
	if !inSlice(outputValues, oFlag) {
		return errors.New(l.Get("Invalid --output parameter value, allowed: table, json, csv"))
	}

	// last
	if lFlag == 0 && dFlag == "" {
		return errors.New(l.Get("Value of one of the parameters should be given: --date or --last"))
	}
	if lFlag < 0 {
		return errors.New(l.Get("Invalid --last parameter value, allowed value > 0"))
	}
	if lFlag > 0 && dFlag != "" {
		return errors.New(l.Get("Only one of the parameters must be given: either --date or --last"))
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
				return errors.New(l.Get("Invalid --date parameter value, allowed values: 'today', 'current', 'YYYY-MM-DD' or 'YYYY-MM-DD: YYYY-MM-DD'"))
			}
		}
	}

	// table
	if cmd == "table" {
		if tFlag == "" {
			return errors.New(l.Get("The --table parameter value is missing, the type of the exchange table should be specified"))
		}
		if !inSlice(tableValues, tFlag) {
			return errors.New(l.Get("Invalid parameter --table value, allowed values: A, B or C"))
		}
	}

	// currency
	var errMessage string

	if cmd == "currency" {
		if cFlag == "" {
			return errors.New(l.Get("No value of parameter --code, currency code should be given"))
		}
		if tFlag == "" {
			return errors.New(l.Get("No value of parameter --table, please specify type of exchange rate table"))
		}

		if !inSlice(tableValues, tFlag) {
			return errors.New(l.Get("Incorrect parameter value --table, allowed values: A, B or C"))
		}

		if tFlag == "A" {
			if !inSlice(currencyValuesA, cFlag) {
				errMessage = l.Get("Incorrect value of the --code parameter, ")
				errMessage += l.Get("valid currency code from those available for Table A is allowed: ")
				errMessage += strings.Join(currencyValuesA, ", ")
				return errors.New(errMessage)
			}
		} else if tFlag == "B" {
			if !inSlice(currencyValuesB, cFlag) {
				errMessage = l.Get("Incorrect value of the --code parameter, ")
				errMessage += l.Get("valid currency code from those available for Table B is allowed: ")
				errMessage += strings.Join(currencyValuesB, ", ")
				return errors.New(errMessage)
			}
		} else if tFlag == "C" {
			if !inSlice(currencyValuesC, cFlag) {
				errMessage = l.Get("Incorrect value of the --code parameter, ")
				errMessage += l.Get("valid currency code from those available for Table C is allowed: ")
				errMessage += strings.Join(currencyValuesC, ", ")
				return errors.New(errMessage)
			}
		}
	}

	return nil
}
