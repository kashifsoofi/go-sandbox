package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("timeLimit", 30, "the time limit for the quiz in seconds")
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

	fmt.Printf("Press any key to start quiz, you have %d seconds to complete.\n", *timeLimit)
	var r rune
	fmt.Scanf("%r", &r)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct := 0

problemloop:
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.question)
		answerCh := make(chan string)

		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			break problemloop
		case answer := <-answerCh:
			if strings.ToLower(strings.TrimSpace(answer)) == p.answer {
				correct++
			}
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
			answer:   strings.ToLower(strings.TrimSpace(r[1])),
		}
	}
	return problems
}
