package main

import (
	"bufio"
	"fmt"
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

func splitAndValidate(l string, n int, s string) ([]string, bool) {
	f := strings.Split(l, s)
	return f, len(f) == n
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

func readThresholds(fs string, s string) (map[string]float64, error) {
	thres := make(map[string]float64)

	f, err := os.Open(fs)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	l := 0
	for sc.Scan() {
		l++
		fields, ok := splitAndValidate(sc.Text(), 2, s)
		if !ok {
			return nil, fmt.Errorf("Thresholds file: format fault at line %d", l)
		}

		val, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			return nil, fmt.Errorf("Thresholds file: cannot convert %s to float64", fields[1])
		}

		key := strings.ToLower(fields[0])
		if prevVal, ok := thres[key]; ok {
			e := fmt.Sprintf("WARN: overwriting previous %s threshold value: %.2f -> %.2f", key, prevVal, val)
			fmt.Fprintln(os.Stderr, e)
		}
		thres[key] = val
	}

	if err := sc.Err(); err != nil {
		return nil, err
	}

	return thres, nil
}

func readMetrics(fs string, s string) ([]Metric, error) {
	var metrics []Metric

	f, err := os.Open(fs)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	l := 0
	for sc.Scan() {
		l++
		fields, ok := splitAndValidate(sc.Text(), 4, s)
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
