package main

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

type metric struct{}

func splitAndValidate(l string, n int, s string) ([]string, bool) {
	f := strings.Split(l, s)
	return f, len(f) == n
}

func readThresholds(fs string, s string) (map[string]float64, error) {
	thres := make(map[string]float64)

	f, err := os.Open(fs)
	if err != nil {
		return thres, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		fields, ok := splitAndValidate(sc.Text(), 2, s)
		if !ok {
			return thres, errors.New("Thresholds file: format fault")
		}

		val, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			return thres, err
		}

		if err := sc.Err(); err != nil {
			return thres, err
		}

		thres[strings.ToLower(fields[0])] = val
	}

	return thres, nil
}

func readMetrics(fs string, s string) (metric, error) {
	m := metric{}

	f, err := os.Open(fs)
	if err != nil {
		return m, err
	}
	defer f.Close()

	return m, nil
}
