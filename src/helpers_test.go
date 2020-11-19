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
