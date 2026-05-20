package main

import (
	"fmt"
	"os"
)

func main() {
	const separator = ","

	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "[-] Must provide exactly 2 argument: [metrics] [thresholds]")
		os.Exit(1)
	}

	metricsFile, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "[-] %s metrics file could not be opened or does not exist\n", os.Args[1])
		os.Exit(1)
	}
	defer metricsFile.Close()

	tFileName := os.Args[2]
	mFileName := os.Args[1]

	fh, err := os.Open(tFileName)
	if err != nil {
		fmt.Errorf("Could not open file %s", tFileName)
		os.Exit(1)
	}

	thresholds, err := readThresholds(fh, separator)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[-] %v\n", err)
	}
	fmt.Println(thresholds)
}
