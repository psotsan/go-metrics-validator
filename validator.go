package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type metric struct {
	Name      string
	Value     float64
	Unit      string
	Timestamp time.Time
}

func splitAndValidate(l string, n int, s string) ([]string, bool) {
	f := strings.Split(l, s)
	return f, len(f) == n
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

		// AÑADIR: Verificar si ya existe el umbral en el mapa antes de agregar
		thres[strings.ToLower(fields[0])] = val
	}

	if err := sc.Err(); err != nil {
		return nil, err
	}

	return thres, nil
}

func readMetrics(fs string, s string) ([]metric, error) {
	var metrics []metric

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
		m := metric{Name: mName, Unit: fields[2]}
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
