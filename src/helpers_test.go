package main

import "testing"

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
