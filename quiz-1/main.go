package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")

	flag.Parse()

	f, err := os.Open(*csvFilename)
	if err != nil {
		fmt.Printf("Error opening CSV file: %s, error: %v", *csvFilename, err)
		os.Exit(1)
	}

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Error parsing CSV file, error: %v", err)
	}

	problems := parseProblems(records)

	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.question)

		var answer string
		fmt.Scanf("%s\n", &answer)
		if answer == p.answer {
			correct++
		}
	}

	fmt.Printf("You have scored %d out of %d.\n", correct, len(problems))
}

type problem struct {
	question string
	answer   string
}

func parseProblems(records [][]string) []problem {
	problems := make([]problem, len(records))
	for i, r := range records {
		problems[i] = problem{
			question: r[0],
			answer:   strings.TrimSpace(r[1]),
		}
	}
	return problems
}
