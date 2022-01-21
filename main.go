package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Config struct {
	PlaybackSpeed     float64 `json:"PlaybackSpeed"`
	TimestampBuffer   int64   `json:"TimestampBuffer"`
	PlaybackStartTime int64   `json:"PlaybackStartTime"`
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
	configPath := "config.json"
	cfgDat, err := os.ReadFile(configPath)
	if err != nil {
		log.Panic()
	}

	var cfg Config
	err = json.Unmarshal(cfgDat, &cfg)
	if err != nil {
		log.Panic()
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
	fmt.Printf("\nStart Time: %s Playback Speed: %.2f Buffer: %d sec\n", convertSeconds(int(cfg.PlaybackStartTime)), cfg.PlaybackSpeed, cfg.TimestampBuffer)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		description, ok := commands[input]
		if ok {
			//Actual time from start
			currentTime := -int(time.Until(now).Seconds())
			//Adjust time passed based on playback speed
			currentTimePlayback := int(float32(currentTime) * float32(cfg.PlaybackSpeed))
			//Adjust time to reflect playback start time and timestamp buffer
			adjustedTime := currentTimePlayback + int(cfg.PlaybackStartTime) - int(cfg.TimestampBuffer)
			timeStamp := convertSeconds(adjustedTime)
			output := fmt.Sprintln(timeStamp, " - ", description)
			p(output)
			f.WriteString(output)
		}

		if input == "exit" {
			return
		}
	}
}
