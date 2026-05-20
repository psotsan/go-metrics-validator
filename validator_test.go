package main

import (
	"testing"
)

func TestSplitAndValidate(t *testing.T) {
	savs := []struct {
		name     string
		line     string
		n        int
		sep      string
		wantStr  []string
		wantBool bool
	}{
		{
			name:     "two fields OK",
			line:     "cpu_usage,65",
			n:        2,
			sep:      ",",
			wantStr:  []string{"cpu_usage", "65"},
			wantBool: true,
		},
		{
			name:     "four fields OK",
			line:     "mem_usage,82.0,%,2025-05-15T10:00:00Z",
			n:        4,
			sep:      ",",
			wantStr:  []string{"mem_usage", "82.0", "%", "2025-05-15T10:00:00Z"},
			wantBool: true,
		},
		{
			name:     "two fields no separator",
			line:     "cpu_usage 65",
			n:        2,
			sep:      ",",
			wantStr:  []string{"cpu_usage 65"},
			wantBool: false,
		},
		{
			name:     "more fields than expected",
			line:     "mem_usage,82.0,%,2025-05-15T10:00:00Z",
			n:        2,
			sep:      ",",
			wantStr:  []string{"mem_usage", "82.0", "%", "2025-05-15T10:00:00Z"},
			wantBool: false,
		},
		{
			name:     "less fields than expected",
			line:     "mem_usage,82.0,",
			n:        4,
			sep:      ",",
			wantStr:  []string{"mem_usage", "82.0"},
			wantBool: false,
		},
		{
			name:     "empty line",
			line:     "",
			n:        2,
			sep:      ",",
			wantStr:  []string{},
			wantBool: false,
		},
		{
			name:     "leading separator",
			line:     ",cpu_usage",
			n:        2,
			sep:      ",",
			wantStr:  []string{"cpu_usage"},
			wantBool: false,
		},
		{
			name:     "trailing separator",
			line:     "cpu_usage,",
			n:        2,
			sep:      ",",
			wantStr:  []string{"cpu_usage"},
			wantBool: false,
		},
		{
			name:     "untrimmed values",
			line:     "    cpu_usage  ,  92 ",
			n:        2,
			sep:      ",",
			wantStr:  []string{"cpu_usage", "92"},
			wantBool: true,
		},
	}

	for _, tt := range savs {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := splitAndValidate(tt.line, tt.n, tt.sep)
			sameStr := true
			if ok == tt.wantBool && len(got) == len(tt.wantStr) {
				for i := range got {
					if got[i] != tt.wantStr[i] {
						sameStr = false
						break
					}
				}
			} else {
				sameStr = false
			}
			if !sameStr {
				t.Errorf("splitAndValidate( %q)\ngot string %v - want %v\ngot ok %t - want %t\n", tt.name, got, tt.wantStr, ok, tt.wantBool)
			}
		})
	}
}
