package main

import (
	"fmt"
	"os"
)

func main() {
	const separator = ","

	if len(os.Args) != 3 {
		fmt.Println("[-] Must provide exactly 2 argument: [metrics] [thresholds]")
		return
	}

	metricsFile, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("[-] %s metrics file could not be opened or does not exist\n", os.Args[1])
		return
	}
	defer metricsFile.Close()

	thresholds, err := readThresholds(os.Args[2], separator)
	if err != nil {
		fmt.Printf("[-] %v\n", err)
		return
	}
	fmt.Println(thresholds)
}
