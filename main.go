package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("quiz.csv")
	var questionsAndAnswers []QuestionAnswer
	if err != nil {
		log.Fatal("Error opening quiz file")
	}
	csvReader := csv.NewReader(file)
	contents, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Error reading contents of quiz file")
	}
	for _, line := range contents {
		answer, err := strconv.Atoi(line[1])
		if err != nil {
			log.Fatal("Unable to properly parse answer in line")
		}
		currentQuestion := QuestionAnswer{Question: line[0], Answer: answer}
		questionsAndAnswers = append(questionsAndAnswers, currentQuestion)
	}
	playGame(questionsAndAnswers)
}

func playGame(questionsAndAnswers []QuestionAnswer) {
	reader := bufio.NewReader(os.Stdin)
	correctAnswers := 0
	fmt.Println("Are you ready to play? Y or N")
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error reading user input from standard input")
	}
	trimResponses(&response)

	if response == "Y" {
		for _, question := range questionsAndAnswers {
			fmt.Printf("What is %s", question.Question)
			response, err = reader.ReadString('\n')
			if err != nil {
				fmt.Println("Unable to read user input, assuming incorrect answer")
			}
			trimResponses(&response)
			answer, err := strconv.Atoi(response)
			if err != nil {
				fmt.Println("You entered an invalid response, assuming incorrect answer")
			}
			if answer == question.Answer {
				fmt.Println("That is correct")
				correctAnswers++
			} else {
				fmt.Println("That is incorrect")
			}
		}
		fmt.Printf("You got %d questions correct", correctAnswers)
	}
}

func trimResponses(response *string) {
	if runtime.GOOS == "windows" {
		*response = strings.TrimRight(*response, "\r\n")
	} else {
		*response = strings.TrimRight(*response, "\n")
	}
}
