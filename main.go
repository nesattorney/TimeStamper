package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func compare(a, b string) bool {
	fmt.Print("\nComparing:\n", a, b, a == b, "\n")
	fmt.Println(a, len(a))
	fmt.Println(b, len(b))
	return a == b
}

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
	p := fmt.Println

	now := time.Now()
	fmt.Println("===Starting Timestamper===")

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		input := scanner.Text()
		comp := "a"
		if input == comp {
			timeStamp := convertSeconds(-int(time.Until(now).Seconds()))
			p(timeStamp, " - ", "sucess")
		}

		if input == "exit" {
			return
		}
	}
}
