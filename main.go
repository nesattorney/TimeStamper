package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
)

type ReplayTimeStamper struct {
	Buffer     int
	Cmds       map[string]string
	OutputFile *os.File
	Exit       chan bool
	TimeStamps map[string][]string
}

//Parses json input, writes timestamp to file if input is registered
func (rt *ReplayTimeStamper) handleInput(w http.ResponseWriter, req *http.Request) {
	var event KeyEvent
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Panic("Invalid input")
	}
	err = json.Unmarshal(data, &event)
	if err != nil {
		log.Panic("json invalid")
	}

	description, ok := rt.Cmds[event.Key]
	var timeStamp string
	if ok {
		seconds := int(event.Timestep) - int(rt.Buffer)
		timeStamp = convertSeconds(seconds)
		output := fmt.Sprint(timeStamp, " - ", description)
		fmt.Println(output)
		rt.OutputFile.WriteString(output + "\n")
	} else if event.Key == "Escape" {
		rt.Exit <- true
		return
	} else {
		fmt.Println("Unbound key:", event.Key)
		return
	}

	currentTimeStamps, ok := rt.TimeStamps[description]
	if !ok {
		currentTimeStamps = make([]string, 0)
	}
	currentTimeStamps = append(currentTimeStamps, timeStamp)
	rt.TimeStamps[description] = currentTimeStamps
}

type KeyEvent struct {
	Title    string  `json:"title"`
	Timestep float64 `json:"timestep"`
	Key      string  `json:"key"`
}

var KeyCmdPathFlag = flag.String("c", "millia.txt", "Key command filename")
var bufferFlag = flag.Int("b", 3, "Timestamp Buffer Window")
var destinationFlag = flag.String("d", "result.txt", "Result destination")

//Converts seconds to hours:minutes:seconds
func convertSeconds(seconds int) string {
	if seconds < 0 {
		seconds = 0
	}

	hours := seconds / 3600
	minutes := seconds / 60
	seconds = seconds - minutes*60 - hours*3600
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

//Formats the output such that timestamps are sorted by type
func prettyPrint(rts ReplayTimeStamper) {
	//Create destination file
	fpretty, err := os.Create(filepath.Join("output", "pretty"+*destinationFlag))
	if err != nil {
		log.Panic()
	}
	defer fpretty.Close()

	for k, v := range rts.TimeStamps {
		fpretty.WriteString("=====" + k + "=====\n")
		for _, line := range v {
			fpretty.WriteString(line + "\n")
		}
		fpretty.WriteString("\n")
	}
}

func main() {
	flag.Parse()
	fmt.Printf("===Using key commands from file: %s ===\n", *KeyCmdPathFlag)
	fmt.Printf("Buffer: %d sec\n", *bufferFlag)
	fmt.Printf("Destination: %s\n", filepath.Join("output", *destinationFlag))

	//Read key commands
	dat, err := os.ReadFile(filepath.Join("keyCommands", *KeyCmdPathFlag))
	if err != nil {
		log.Panic("Invalid command path")
	}
	keyCommands := map[string]string{}

	//Add them to map with Key=keybind, Value=TimestampDescription
	inputLines := strings.Split(string(dat), "\n")
	for _, line := range inputLines {
		line = strings.TrimSuffix(line, "\n")
		line = strings.TrimSuffix(line, "\r")
		command := strings.Split(line, "=")
		if len(command) != 2 {
			log.Panic("Invalid input:", command)
		}
		keyCommands[command[0]] = command[1]
	}

	//Create destination file
	f, err := os.Create(filepath.Join("output", *destinationFlag))
	if err != nil {
		log.Panic()
	}
	defer f.Close()

	//List available commands to client
	fmt.Println("KEY COMMANDS:")
	for k, v := range keyCommands {
		fmt.Println(k, ":", v)
	}
	fmt.Println("Escape", ":", "exits program")
	fmt.Println("\n===Starting TimeStamper===")

	rts := ReplayTimeStamper{
		Buffer:     int(*bufferFlag),
		Cmds:       keyCommands,
		OutputFile: f,
		Exit:       make(chan bool),
		TimeStamps: make(map[string][]string),
	}
	http.HandleFunc("/", rts.handleInput)
	go http.ListenAndServe("localhost:9393", nil)

	// Exits when rts receives the Escape key as input
	// Or when process is interrupted with ctrl+c
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	select {
	case <-rts.Exit:
	case <-sig:
	}

	//sort timestamps and write to file in a cleaner format
	prettyPrint(rts)
	fmt.Println("Shutting down...")
}
