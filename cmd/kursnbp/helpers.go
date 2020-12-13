// auxiliary functions
package main

import (
	"errors"
	"regexp"
	"strings"

	"github.com/pjaskulski/nbpapi"
)

var outputValues = []string{"table", "json", "csv", "xml"}

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
func checkArg(cmd, tFlag, dFlag string, lFlag int, oFlag, cFlag string) error {

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
			re10 := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
			isValid = re10.MatchString(dFlag)
		} else if len(dFlag) == 21 {
			re21 := regexp.MustCompile(`\d{4}-\d{2}-\d{2}\:\d{4}-\d{2}-\d{2}`)
			isValid = re21.MatchString(dFlag)
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
	if !inSlice(nbpapi.TableValues, tFlag) {
		return errors.New("Invalid parameter --table value, allowed values: A, B or C")
	}

	return nil
}

// currency code check
func checkArgCurrency(tFlag, cFlag string) error {
	var errMessage string

	if cFlag == "" {
		return errors.New("No value of parameter --code, currency code should be given")
	}
	if tFlag == "" {
		return errors.New("No value of parameter --table, please specify type of exchange rate table")
	}

	if !inSlice(nbpapi.TableValues, tFlag) {
		return errors.New("Incorrect parameter value --table, allowed values: A, B or C")
	}

	if tFlag == "A" {
		if !inSlice(nbpapi.CurrencyValuesA, cFlag) {
			errMessage = "Incorrect value of the --code parameter, "
			errMessage += "valid currency code from those available for Table A is allowed: "
			errMessage += strings.Join(nbpapi.CurrencyValuesA, ", ")
			return errors.New(errMessage)
		}
	} else if tFlag == "B" {
		if !inSlice(nbpapi.CurrencyValuesB, cFlag) {
			errMessage = "Incorrect value of the --code parameter, "
			errMessage += "valid currency code from those available for Table B is allowed: "
			errMessage += strings.Join(nbpapi.CurrencyValuesB, ", ")
			return errors.New(errMessage)
		}
	} else if tFlag == "C" {
		if !inSlice(nbpapi.CurrencyValuesC, cFlag) {
			errMessage = "Incorrect value of the --code parameter, "
			errMessage += "valid currency code from those available for Table C is allowed: "
			errMessage += strings.Join(nbpapi.CurrencyValuesC, ", ")
			return errors.New(errMessage)
		}
	}

	return nil
}
