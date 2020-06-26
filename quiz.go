package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

var csvFilepath string
var timeout int
var start time.Time = time.Now()
var csvLines [][]string
var result quizResult = quizResult{0, 0}

type problemsData struct {
	Question string
	Result   string
}

type quizResult struct {
	Corrects int
	Total    int
}

func processFileArguments() {
	flag.StringVar(&csvFilepath, "csvfile", "problems.csv", "path to csv with the problems")
	flag.IntVar(&timeout, "timeout", 30, "quiz timeout in seconds")
	flag.Parse()
}

func countSeconds(seconds float64) {
	for range time.Tick(1 * time.Second) {
		if time.Now().Sub(start).Seconds() >= seconds {
			fmt.Println("Time's up")
			printResults()
			os.Exit(1)
		}
	}
}

func askQuestion(question string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(question)
	text, _ := reader.ReadString('\n')
	return strings.Trim(text, "\n")
}

func isCorrectAnswer(answer string, expectedAnswer string) bool {
	return strings.Compare(answer, expectedAnswer) == 0
}

func printResults() {
	fmt.Println("Number of questions answered correctly:")
	fmt.Println(result.Corrects)
	fmt.Println("Total number of questions")
	fmt.Println(result.Total)
}

func main() {

	processFileArguments()
	askQuestion(fmt.Sprintf("You have %d second(s) to finish the quiz\nPlease press enter do start it:", timeout))
	go countSeconds(float64(timeout))

	csvFile, err := os.Open(csvFilepath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
		return
	}

	result.Total = len(csvLines)

	for _, line := range csvLines {
		problem := problemsData{
			Question: line[0],
			Result:   strings.Trim(line[1], "\n"),
		}
		answer := askQuestion(problem.Question)
		isCorrect := isCorrectAnswer(answer, problem.Result)
		if isCorrect {
			result.Corrects++
		}
	}
	fmt.Println("No more questions left")
	printResults()
}
