package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type ReplayTimeStamper struct {
	Buffer     int
	Cmds       map[string]string
	OutputFile *os.File
	Exit       chan bool
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

	if ok {
		seconds := int(event.Timestep) - int(rt.Buffer)
		timeStamp := convertSeconds(seconds)
		output := fmt.Sprint(timeStamp, " - ", description)
		fmt.Println(output)
		rt.OutputFile.WriteString(output + "\n")
	} else if event.Key == "Escape" {
		rt.Exit <- true
	} else {
		fmt.Println("Unbound key:", event.Key)
	}
}

func main() {
	flag.Parse()
	fmt.Printf("===Using key commands from file: %s ===\n", *KeyCmdPathFlag)
	fmt.Printf("Buffer: %d sec\n", *bufferFlag)

	//Read key commands
	dat, err := os.ReadFile(filepath.Join("keyCommands", *KeyCmdPathFlag))
	if err != nil {
		log.Panic("Invalid command path")
	}
	keyCommands := map[string]string{}

	//Add them to map with Key=keybind, Value=TimestampDescription
	inputLines := strings.Split(string(dat), "\n")
	for _, line := range inputLines {
		command := strings.Split(line, "=")
		if len(command) != 2 {
			log.Panic("Invalid input:", command)
		}
		keyCommands[command[0]] = command[1]
	}

	//Create destination file
	f, err := os.Create(*destinationFlag)
	if err != nil {
		log.Panic()
	}
	defer f.Close()

	//List available commands to client
	fmt.Println("KEY COMMANDS:")
	for k, v := range keyCommands {
		fmt.Println(k, ":", v)
	}
	fmt.Println("exit", ":", "exits program")
	fmt.Println("\n===Starting TimeStamper===")

	rts := ReplayTimeStamper{
		Buffer:     int(*bufferFlag),
		Cmds:       keyCommands,
		OutputFile: f,
		Exit:       make(chan bool),
	}
	http.HandleFunc("/", rts.handleInput)
	go http.ListenAndServe("localhost:9393", nil)

	//Exits when rts receives the Escape key as input
	<-rts.Exit
	fmt.Println("Shutting down...")
}
