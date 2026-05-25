package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type Metric struct {
	Name      string
	Value     float64
	Unit      string
	Timestamp time.Time
}

type Threshold struct {
	value        float64
	isUpperLimit bool
}

func splitAndValidate(l string, n int) ([]string, bool) {
	var ret []string
	f := strings.Split(l, separator)
	for _, st := range f {
		st = strings.Trim(st, " ")
		if st != "" {
			ret = append(ret, st)
		}
	}
	return ret, len(ret) == n
}

func readThresholds(r io.Reader) (thresholds map[string]Threshold, err error) {
	thresholds = make(map[string]Threshold)

	sc := bufio.NewScanner(r)
	l := 0
	for sc.Scan() {
		l++
		line := sc.Text()

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields, ok := splitAndValidate(line, 3)
		if !ok {
			return nil, fmt.Errorf("Thresholds file: format fault at line %d", l)
		}

		val, e := strconv.ParseFloat(fields[1], 64)
		if e != nil {
			return nil, fmt.Errorf("Thresholds file: could not convert %s to float64 at line %d", fields[1], l)
		}

		limitType := strings.ToLower(fields[2])
		if limitType != "max" && limitType != "min" {
			return nil, fmt.Errorf("Thresholds file: unrecognized type of limit at line %d", l)
		}

		key := strings.ToLower(fields[0])
		if prevThres, ok := thresholds[key]; ok {
			e := fmt.Sprintf("WARN: overwriting previous %s threshold value: %.2f -> %.2f", key, prevThres.value, val)
			fmt.Fprintln(os.Stderr, e)
		}

		isUpperLimit := limitType == "max"
		thresholds[key] = Threshold{value: val, isUpperLimit: isUpperLimit}
	}

	if err := sc.Err(); err != nil {
		return nil, err
	}
	return
}

func readMetrics(r io.Reader) ([]Metric, error) {
	var metrics []Metric

	sc := bufio.NewScanner(r)
	l := 0
	for sc.Scan() {
		l++
		line := sc.Text()

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields, ok := splitAndValidate(line, 4)
		if !ok {
			e := fmt.Sprintf("Metrics file: format fault at line %d", l)
			fmt.Fprintln(os.Stderr, e)
			continue
		}

		mName := strings.ToLower(fields[0])
		m := Metric{Name: mName, Unit: fields[2]}
		if v, err := strconv.ParseFloat(fields[1], 64); err == nil {
			m.Value = v
		} else {
			e := fmt.Sprintf("Metrics file, line %d: cannot convert %s to float 64", l, fields[1])
			fmt.Fprintln(os.Stderr, e)
			continue
		}

		if v, err := time.Parse("2006-01-02T15:04:05Z", fields[3]); err == nil {
			m.Timestamp = v
		} else {
			e := fmt.Sprintf("Metrics file, line %d: cannot convert %s to time.Time", l, fields[3])
			fmt.Fprintln(os.Stderr, e)
			continue
		}

		metrics = append(metrics, m)
	}

	if err := sc.Err(); err != nil {
		return metrics, err
	}
	return metrics, nil
}

func evaluateMetrics(metrics []Metric, thresholds map[string]Threshold) (warnMetrics []Metric) {
	warnMetrics = make([]Metric, 0)

	for _, m := range metrics {
		var t Threshold
		var ok bool

		if t, ok = thresholds[m.Name]; !ok {
			continue
		}

		if t.isUpperLimit && m.Value > t.value {
			warnMetrics = append(warnMetrics, m)
			continue
		}

		if !t.isUpperLimit && m.Value < t.value {
			warnMetrics = append(warnMetrics, m)
			continue
		}
	}
	return warnMetrics
}
