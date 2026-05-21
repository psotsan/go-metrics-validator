package main

import (
	"strings"
	"testing"
	"time"
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
				t.Fatalf("splitAndValidate(%q)\ngot string %v - want %v\ngot ok %t - want %t\n", tt.name, got, tt.wantStr, ok, tt.wantBool)
			}
		})
	}
}

func TestReadThresholds(t *testing.T) {
	rts := []struct {
		name    string
		input   string
		want    map[string]Threshold
		wantErr bool
	}{
		{
			name:  "correct single line",
			input: "cpu_usage,80,max",
			want:  map[string]Threshold{"cpu_usage": {value: 80.0, isUpperLimit: true}},
		},
		{
			name:  "correct multiple lines",
			input: "cpu_usage,80,max\nmem_usage,90,max\ndisk_free,10,min",
			want: map[string]Threshold{
				"cpu_usage": {value: 80.0, isUpperLimit: true},
				"mem_usage": {value: 90.0, isUpperLimit: true},
				"disk_free": {value: 10.0, isUpperLimit: false},
			},
		},
		{
			name:  "spaced line",
			input: " cpu_usage , 80 , max ",
			want:  map[string]Threshold{"cpu_usage": {value: 80.0, isUpperLimit: true}},
		},
		{
			name:  "empty line",
			input: "",
			want:  map[string]Threshold{},
		},
		{
			name:  "commented line",
			input: "# This is a comment",
			want:  map[string]Threshold{},
		},
		{
			name:  "commented line and valid line",
			input: "# This is a comment\nmem_usage,85,max",
			want:  map[string]Threshold{"mem_usage": {value: 85.0, isUpperLimit: true}},
		},
		{
			name:  "uppercase line",
			input: "CPU_USAGE,80,MAX",
			want:  map[string]Threshold{"cpu_usage": {value: 80.0, isUpperLimit: true}},
		},
		{
			name:  "overwritten metrics. Should show a warning",
			input: "cpu_usage,80,max\ncpu_usage,85,max",
			want:  map[string]Threshold{"cpu_usage": {value: 85.0, isUpperLimit: true}},
		},
		{
			name:    "line with less than 3 fields",
			input:   "cpu_usage, 80",
			wantErr: true,
		},
		{
			name:    "line with more than 3 fields",
			input:   "cpu_usage, 80, max, extra",
			wantErr: true,
		},
		{
			name:    "non numeric second field",
			input:   "cpu_usage, eighty, max",
			wantErr: true,
		},
		{
			name:    "third field not max or min",
			input:   "cpu_usage, 80, upper",
			wantErr: true,
		},
	}

	for _, tt := range rts {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			thres, err := readThresholds(r, ",")

			if tt.wantErr && err == nil {
				t.Fatalf("readThresholds (%q) - error not present when expected", tt.name)
			}

			if len(thres) == len(tt.want) {
				for k, v := range tt.want {
					if thres[k] != v {
						t.Errorf("readThresholds (%q) - expected %s = %v - got %s = %v", tt.name, k, v, k, thres[k])
					}
				}
			}
		})
	}
}

func TestReadMetrics(t *testing.T) {
	rms := []struct {
		name  string
		input string
		want  []Metric
		// wantErr bool
	}{
		{
			name:  "correct single line",
			input: "cpu_usage,75.2,%,2025-05-15T10:00:00Z",
			want: []Metric{
				{
					Name:      "cpu_usage",
					Value:     75.2,
					Unit:      "%",
					Timestamp: time.Date(2025, 5, 15, 10, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name:  "correct multiple lines",
			input: "cpu_usage,75.2,%,2025-05-15T10:00:00Z\ndisk_free,15.3,GB,2025-05-15T10:00:00Z",
			want: []Metric{
				{
					Name:      "cpu_usage",
					Value:     75.2,
					Unit:      "%",
					Timestamp: time.Date(2025, 5, 15, 10, 0, 0, 0, time.UTC),
				},
				{
					Name:      "disk_free",
					Value:     15.3,
					Unit:      "GB",
					Timestamp: time.Date(2025, 5, 15, 10, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name:  "spaced line",
			input: " cpu_usage, 75.2, %, 2025-05-15T10:00:00Z ",
			want: []Metric{
				{
					Name:      "cpu_usage",
					Value:     75.2,
					Unit:      "%",
					Timestamp: time.Date(2025, 5, 15, 10, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name:  "empty line",
			input: "",
			want:  []Metric{},
		},
		{
			name:  "commented line",
			input: "# This is a comment",
			want:  []Metric{},
		},
		{
			name:  "comented line and valid line",
			input: "# This is a comment\ncpu_usage,75.2,%,2025-05-15T10:00:00Z",
			want: []Metric{
				{
					Name:      "cpu_usage",
					Value:     75.2,
					Unit:      "%",
					Timestamp: time.Date(2025, 5, 15, 10, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name:  "uppercase line",
			input: "CPU_USAGE,75.2,%,2025-05-15T10:00:00Z",
			want: []Metric{
				{
					Name:      "cpu_usage",
					Value:     75.2,
					Unit:      "%",
					Timestamp: time.Date(2025, 5, 15, 10, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name:  "line with less than 4 fields. Should show warning",
			input: "CPU_USAGE,75.2,%",
			want:  []Metric{},
		},
		{
			name:  "line with more than 4 fields. Should show warning",
			input: "CPU_USAGE,75.2,%,2025-05-15T10:00:00Z,extra",
			want:  []Metric{},
		},
		{
			name:  "non-numeric second field. Should show warning",
			input: "cpu_usage,seventy five,%,2025-05-15T10:00:00Z",
			want:  []Metric{},
		},
	}
}
