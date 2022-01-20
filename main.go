package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func convertSeconds(seconds int) string {
	if seconds < 0 {
		seconds = 0
	}

	hours := seconds / 3600
	minutes := seconds / 60
	seconds = seconds - minutes*60 - hours*3600
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func main() {
	inputPath := "input.txt"
	dat, err := os.ReadFile(inputPath)
	if err != nil {
		log.Panic()
	}
	commands := map[string]string{}
	inputLines := strings.Split(string(dat), "\n")
	for _, line := range inputLines {
		command := strings.Split(line, "=")
		if len(command) != 2 {
			log.Panic("Invalid input:", command)
		}
		commands[command[0]] = command[1]
	}

	p := fmt.Print

	f, err := os.Create("result.txt")
	if err != nil {
		log.Panic()
	}

	now := time.Now()
	p("COMMANDS:\n")
	for k, v := range commands {
		fmt.Println(k, ":", v)
	}
	fmt.Println("exit", ":", "exits program")

	p("\n===Starting TimeStamper===\n")
	offset := 2
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		description, ok := commands[input]
		if ok {
			timeStamp := convertSeconds(-(int(time.Until(now).Seconds()) + offset))
			output := fmt.Sprintln(timeStamp, " - ", description)
			p(output)
			f.WriteString(output)
		}

		if input == "exit" {
			return
		}
	}
}
