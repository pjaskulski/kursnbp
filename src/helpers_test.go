package main

import "testing"

func TestInSlice(t *testing.T) {

	tests := []struct {
		name  string
		slice []string
		text  string
		want  bool
	}{
		{
			name:  "Poprawny kod",
			slice: []string{"CHF", "EUR", "HKD", "USD", "DKK", "GBP"},
			text:  "CHF",
			want:  true,
		},
		{
			name:  "Niepoprawny kod",
			slice: []string{"CHF", "EUR", "HKD", "USD", "DKK", "GBP"},
			text:  "JPY",
			want:  false,
		},
		{
			name:  "Niepoprawny typ tabeli",
			slice: []string{"A", "B", "C"},
			text:  "E",
			want:  false,
		},
		{
			name:  "Poprawny typ tabeli",
			slice: []string{"A", "B", "C"},
			text:  "C",
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := inSlice(tt.slice, tt.text)
			if result != tt.want {
				t.Errorf("oczekiwano: %t; otrzymano: %t", tt.want, result)
			}
		})
	}
}

// checkArg(cmd string, tFlag string, dFlag string, lFlag int, oFlag string, cFlag string)

func TestCheckArg(t *testing.T) {
	tests := []struct {
		name        string
		cmd         string
		tFlag       string
		dFlag       string
		lFlag       int
		oFlag       string
		cFlag       string
		occursError bool
	}{
		{
			name:        "CmdCurrencyShouldBeNoError",
			cmd:         "currency",
			tFlag:       "A",
			dFlag:       "2020-11-12",
			lFlag:       0,
			oFlag:       "table",
			cFlag:       "CHF",
			occursError: false,
		},
		{
			name:        "CmdCurrencyShouldBeErrorComesFromDate",
			cmd:         "currency",
			tFlag:       "A",
			dFlag:       "2020-MM-12",
			lFlag:       0,
			oFlag:       "table",
			cFlag:       "CHF",
			occursError: true,
		},
		{
			name:        "CmdTableShouldBeNoError",
			cmd:         "table",
			tFlag:       "A",
			dFlag:       "2020-11-12",
			lFlag:       0,
			oFlag:       "table",
			cFlag:       "",
			occursError: false,
		},
		{
			name:        "CmdTableShouldBeErrorComesFromType",
			cmd:         "table",
			tFlag:       "D",
			dFlag:       "2020-11-12",
			lFlag:       0,
			oFlag:       "table",
			cFlag:       "",
			occursError: true,
		},
		{
			name:        "CmdGoldShouldBeNoError",
			cmd:         "gold",
			tFlag:       "",
			dFlag:       "2020-11-19",
			lFlag:       0,
			oFlag:       "table",
			cFlag:       "",
			occursError: false,
		},
		{
			name:        "CmdGoldShouldBeErrorComesFromDateAndLastAtTheSameTime",
			cmd:         "gold",
			tFlag:       "",
			dFlag:       "2020-11-19",
			lFlag:       10,
			oFlag:       "table",
			cFlag:       "",
			occursError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkArg(tt.cmd, tt.tFlag, tt.dFlag, tt.lFlag, tt.oFlag, tt.cFlag)
			if tt.occursError == false && result != nil {
				t.Errorf("expected: no errors; received: error")
			} else if tt.occursError == true && result == nil {
				t.Errorf("expected: error; received: nil")
			}
		})
	}
}
