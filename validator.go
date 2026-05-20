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

type threshold struct {
	value        float64
	isUpperLimit bool
}

func splitAndValidate(l string, n int, s string) ([]string, bool) {
	var ret []string
	f := strings.Split(l, s)
	for _, st := range f {
		st = strings.Trim(st, " ")
		if st != "" {
			ret = append(ret, st)
		}
	}
	return ret, len(ret) == n
}

func validateMetric(m Metric, t map[string]float64) error {
	if _, ok := t[m.Name]; !ok {
		return fmt.Errorf("Metric %s not found in thresholds file", m.Name)
	}
	return nil
}

func checkMetric(m Metric, t map[string]float64, upperThres bool) (exceedsThres bool, err error) {
	if e := validateMetric(m, t); e != nil {
		return false, e
	}

	if (upperThres && m.Value > t[m.Name]) || (!upperThres && m.Value < t[m.Name]) {
		return true, nil
	}
	return false, nil
}

func readThresholds(r io.Reader, sep string) (thresholds map[string]threshold, err error) {
	thresholds = make(map[string]threshold)

	sc := bufio.NewScanner(r)
	l := 0
	for sc.Scan() {
		l++
		line := sc.Text()

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields, ok := splitAndValidate(line, 3, sep)
		if !ok {
			return nil, fmt.Errorf("Thresholds file: format fault at line %d", l)
		}

		val, e := strconv.ParseFloat(fields[1], 64)
		if e != nil {
			return nil, fmt.Errorf("Thresholds file: format fault at line %d", l)
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
		thresholds[key] = threshold{value: val, isUpperLimit: isUpperLimit}
	}

	if err := sc.Err(); err != nil {
		return nil, err
	}

	return
}

func readMetrics(r io.Reader, sep string) ([]Metric, error) {
	var metrics []Metric

	sc := bufio.NewScanner(r)
	l := 0
	for sc.Scan() {
		l++
		fields, ok := splitAndValidate(sc.Text(), 4, sep)
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

func evaluateMetrics(m []Metric, t map[string]float64) (warnings int) {
	return warnings
}
