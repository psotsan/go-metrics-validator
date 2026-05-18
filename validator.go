package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func getFields(l string, n int, s string) ([]string, bool) {
	f := strings.Split(l, s)
	return f, len(f) == n
}

func readThresholds(fs string, s string) (map[string]int, error) {
	thres := make(map[string]int)

	f, err := os.Open(fs)
	if err != nil {
		return thres, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		fields, ok := getFields(sc.Text(), 2, s)
		if !ok {
			return thres, errors.New("Thresholds file: format fault")
		}
		fmt.Println(fields)
	}

	return thres, nil
}
