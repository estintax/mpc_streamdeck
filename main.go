package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/andreykaipov/goobs"
	"github.com/tarm/serial"
)

var client *goobs.Client

func enterSerial(prompt string, exit bool) string {
	for {
		fmt.Print(prompt)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()
		if input != "" {
			return input
		}
		if exit {
			break
		}
	}
	return ""
}

func reinitAnswer() {
	answer := enterSerial("Reinit config file? [Yes/No] ", false)
	if strings.ToLower(answer) == "yes" {
		if resetConfig() {
			enterSerial("Restart me now. Enter to exit.", true)
			os.Exit(1)
			return
		}
	} else {
		enterSerial("Then restart the program. Maybe this will help.", true)
		os.Exit(1)
		return
	}

}

func main() {
	fmt.Println("MPC Streamdeck (C) 2024 Maksim Pinigin")
	if !initConfig() {
		config.serialPort = enterSerial("Enter serial port: ", false)
		config.host = enterSerial("Enter OBS API host: ", false)
		config.password = enterSerial("Enter OBS API password: ", false)
		config.dlScriptPath = enterSerial("Enter path to the your DinoLang script: ", false)
		saveConfig()
	}
	initDL()
	var err error
	client, err = goobs.New(config.host, goobs.WithPassword(config.password))
	if err != nil {
		log.Println("Unable to create OBS API client:", err.Error())
		reinitAnswer()
		return
	}
	defer client.Disconnect()
	version, err := client.General.GetVersion()
	if err != nil {
		log.Println("Unable to get OBS version:", err.Error())
		reinitAnswer()
		return
	}
	log.Println("OBS version:", version.ObsVersion)
	c := &serial.Config{Name: config.serialPort, Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Println("Unable to init serial:", err.Error())
		reinitAnswer()
	}
	defer s.Close()
	reader := bufio.NewReader(s)
	log.Println("Serial port opened. Working.")
	runScript(-1)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err.Error())
			time.Sleep(1 * time.Second)
		} else {
			input = strings.Trim(input, "\r\n")
			log.Println("Key:", input)
			key, _ := strconv.Atoi(input)
			if key >= 1 && key <= 4 {
				runScript(key)
			} else {
				log.Println("Bad key pressed.")
			}
		}
	}
}
