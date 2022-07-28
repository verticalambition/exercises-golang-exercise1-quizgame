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
	"time"
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
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Are you ready to play? Y or N")
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error reading user input from standard input")
	}
	trimResponses(&response)
	if response != "Y" {
		fmt.Println("Exiting application")
		os.Exit(0)
	}
	answerChannel := make(chan int)
	quitChannel := make(chan string)
	finalScore := 0
	go gameTimer(10, quitChannel)
	go playGame(questionsAndAnswers, reader, answerChannel, quitChannel)
readChannel:
	for {
		select {
		case <-quitChannel:
			fmt.Printf("You got %d questions correct", finalScore)
			break readChannel
		case <-answerChannel:
			finalScore++
		}
	}
}

func playGame(questionsAndAnswers []QuestionAnswer, reader *bufio.Reader, answerChannel chan int, quitChannel chan string) {

	for _, question := range questionsAndAnswers {
		fmt.Printf("What is %s\n", question.Question)
		response, err := reader.ReadString('\n')
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
			answerChannel <- 1
		} else {
			fmt.Println("That is incorrect")
		}
	}
	quitChannel <- "Quiz Completed"
}
func trimResponses(response *string) {
	if runtime.GOOS == "windows" {
		*response = strings.TrimRight(*response, "\r\n")
	} else {
		*response = strings.TrimRight(*response, "\n")
	}
}

func gameTimer(timeLimit int, quitChannel chan string) {
	time.AfterFunc(time.Second*time.Duration(timeLimit), func() {
		quitChannel <- "Time is Up"
	})
}
