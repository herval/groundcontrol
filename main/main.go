package main

import (
	"flag"
	"log"
	"fmt"
	"syscall"
	"os"
	"bufio"
	"github.com/herval/groundcontrol"
	"strings"
	"strconv"
	"errors"
)

func main() {
	pipe := "groundcontrol.lock"

	port := flag.String("port", "", "/dev/ port connected to your Ground Control device")
	mode := flag.String("mode", "daemon", "Run the daemon or send commands to it (daemon | cmd)")
	device := flag.String("device", "", "The output to send a command to (led <0-2> | display)")
	state := flag.String("state", "", "The new state of the output (on/off for leds, text for display")
	//poll := flag.String("poll", "", "poll for changes on state (any change will return a JSON response)")

	flag.Parse()

	// run as a PID
	if *mode == "daemon" {
		if *port == "" {
			flag.Usage()
			log.Fatal("Missing params")
		}

		fmt.Println("Connecting to Ground Control on port", *port)
		control := groundcontrol.NewGroundControl(*port)
		//control.Connect()

		reader := setupPipe(pipe)
		for {
			device, err := reader.ReadBytes('\n')
			state, err2 := reader.ReadBytes('\n')
			if err != nil || err2 != nil {
				fmt.Print("Couldn't load command: ", err.Error())
			} else {
				funct, err := parseCommand(string(device), string(state), control)
				if err != nil {
					fmt.Println("Couldn't parse command: ", err.Error())
				} else {
					funct()
				}
			}
		}

	} else { // read commands sent to the PID and respond with JSON

		if *device == "" || *state == "" {
			log.Fatal("Missing params")
		}

		file, err := os.OpenFile(pipe, os.O_RDWR | os.O_APPEND, os.ModeNamedPipe)
		handle(err)

		file.WriteString(fmt.Sprintf(
			"%s %s\n", *device, *state,
		))

		fmt.Println("Done!")
	}
}

func setupPipe(filename string) *bufio.Reader {
	os.Remove(filename)

	handle(syscall.Mkfifo(filename, 0666))

	file, err := os.OpenFile(filename, os.O_CREATE, os.ModeNamedPipe)
	handle(err)

	reader := bufio.NewReader(file)

	return reader
}

func parseCommand(device string, state string, control *groundcontrol.GroundControl) (callback func(), err error) {

	switch {
	case strings.HasPrefix(device, "led"):
		port := strings.Split(device, " ")[1]
		return func() {
			if state == "on" {
				control.Leds[port].On()
			} else {
				control.Leds[port].Off()
			}
		}, nil

	case device == "display":
		return func() {
			control.Display.Write(state)
		}, nil
	}

	return nil, errors.New("Couldn't parse param")
}

func intValue(input string) int {
	val, _ := strconv.Atoi(input)
	return val
}

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
