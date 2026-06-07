package main

import (
	"fmt"
	"os"
)

const separator = ","

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "[-] Must provide exactly 2 argument: [metrics] [thresholds]")
		os.Exit(1)
	}

	tFileName := os.Args[2]
	mFileName := os.Args[1]

	r, err := os.Open(tFileName)
	if err != nil {
		e := fmt.Errorf("could not open file %s", tFileName)
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
	defer func() {
		if err := r.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "warning: error closing thresholds file: %v\n", err)
		}
	}()

	thresholds, err := readThresholds(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[-] %v\n", err)
		os.Exit(1)
	}

	r, err = os.Open(mFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[-] %s metrics file could not be opened or does not exist\n", os.Args[1])
		os.Exit(1)
	}
	defer func() {
		if err := r.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "warning: error closing thresholds file: %v\n", err)
		}
	}()

	metrics, err := readMetrics(r)
	metrics = evaluateMetrics(metrics, thresholds)

	for _, m := range metrics {
		fmt.Printf("%s = %.2f (threshold %.2f) [WARNING]\n", m.Name, m.Value, thresholds[m.Name].value)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
