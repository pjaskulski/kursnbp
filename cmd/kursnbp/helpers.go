// auxiliary functions
package main

import (
	"errors"
	"regexp"
	"strings"
)

var tableValues = []string{"A", "B", "C"}
var outputValues = []string{"table", "json", "csv", "xml"}
var currencyValuesA = []string{"THB", "USD", "AUD", "HKD", "CAD", "NZD", "SGD", "EUR", "HUF", "CHF",
	"GBP", "UAH", "JPY", "CZK", "DKK", "ISK", "NOK", "SEK", "HRK", "RON",
	"BGN", "TRY", "ILS", "CLP", "PHP", "MXN", "ZAR", "BRL", "MYR", "RUB",
	"IDR", "INR", "KRW", "CNY", "XDR"}
var currencyValuesB = []string{"MGA", "PAB", "ETB", "AFN", "VES", "BOB", "CRC", "SVC", "NIO", "GMD",
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
var currencyValuesC = []string{"USD", "AUD", "CAD", "EUR", "HUF", "CHF", "GBP", "JPY", "CZK", "DKK", "NOK",
	"SEK", "XDR"}

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

	// output
	if oFlag == "" {
		return errors.New("No --output parameter value, output format must be specified")
	} else if !inSlice(outputValues, oFlag) {
		return errors.New("Invalid --output parameter value, allowed: table, json, csv, xml")
	}

	// last
	if lFlag == 0 && dFlag == "" {
		return errors.New("Value of one of the parameters should be given: --date or --last")
	} else if lFlag != 0 && dFlag != "" {
		return errors.New("Only one of the parameters must be given: either --date or --last")
	} else if lFlag < 0 {
		return errors.New("Invalid --last parameter value, allowed value > 0")
	}

	// date
	err := chkArgDate(dFlag, lFlag)
	if err != nil {
		return err
	}

	// table or currency
	switch cmd {
	case "table":
		return checkArgTable(tFlag)
	case "currency":
		return checkArgCurrency(tFlag, cFlag)
	}

	return nil
}

// check date or last
func chkArgDate(dFlag string, lFlag int) error {
	var isValid bool = true

	if lFlag == 0 && dFlag != "" && dFlag != "today" && dFlag != "current" {
		if len(dFlag) != 10 && len(dFlag) != 21 {
			isValid = false
		} else if len(dFlag) == 10 {
			re10 := regexp.MustCompile("\\d{4}-\\d{2}-\\d{2}")
			isValid = re10.MatchString(dFlag) == true
		} else if len(dFlag) == 21 {
			re21 := regexp.MustCompile("\\d{4}-\\d{2}-\\d{2}\\:\\d{4}-\\d{2}-\\d{2}")
			isValid = re21.MatchString(dFlag) == true
		}
		if !isValid {
			return errors.New("Invalid --date parameter value, allowed values: 'today', 'current', 'YYYY-MM-DD' or 'YYYY-MM-DD:YYYY-MM-DD'")
		}
	}

	return nil
}

// table type check
func checkArgTable(tFlag string) error {
	if tFlag == "" {
		return errors.New("The --table parameter value is missing, the type of the exchange table should be specified")
	}
	if !inSlice(tableValues, tFlag) {
		return errors.New("Invalid parameter --table value, allowed values: A, B or C")
	}

	return nil
}

// currency code check
func checkArgCurrency(tFlag string, cFlag string) error {
	var errMessage string

	if cFlag == "" {
		return errors.New("No value of parameter --code, currency code should be given")
	}
	if tFlag == "" {
		return errors.New("No value of parameter --table, please specify type of exchange rate table")
	}

	if !inSlice(tableValues, tFlag) {
		return errors.New("Incorrect parameter value --table, allowed values: A, B or C")
	}

	if tFlag == "A" {
		if !inSlice(currencyValuesA, cFlag) {
			errMessage = "Incorrect value of the --code parameter, "
			errMessage += "valid currency code from those available for Table A is allowed: "
			errMessage += strings.Join(currencyValuesA, ", ")
			return errors.New(errMessage)
		}
	} else if tFlag == "B" {
		if !inSlice(currencyValuesB, cFlag) {
			errMessage = "Incorrect value of the --code parameter, "
			errMessage += "valid currency code from those available for Table B is allowed: "
			errMessage += strings.Join(currencyValuesB, ", ")
			return errors.New(errMessage)
		}
	} else if tFlag == "C" {
		if !inSlice(currencyValuesC, cFlag) {
			errMessage = "Incorrect value of the --code parameter, "
			errMessage += "valid currency code from those available for Table C is allowed: "
			errMessage += strings.Join(currencyValuesC, ", ")
			return errors.New(errMessage)
		}
	}

	return nil
}
