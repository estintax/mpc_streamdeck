package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
)

type configType struct {
	host         string
	password     string
	serialPort   string
	dlScriptPath string
}

var config configType = configType{host: "", password: "", serialPort: ""}

func getConfigPath() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("UserProfile") + "/mpc_streamdeck.conf"
	} else if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		return os.Getenv("HOME") + "/mpc_streamdeck.conf"
	} else {
		log.Fatalln("Unsupported OS")
		return ""
	}
}

func saveConfig() bool {
	file, err := os.OpenFile(getConfigPath(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatalln("Unable to create config file:", err.Error())
		return false
	}
	defer file.Close()
	file.Write([]byte(fmt.Sprintf("%s;%s;%s;%s;", config.host, config.password, config.serialPort, config.dlScriptPath)))
	return true
}

func initConfig() bool {
	file, err := os.Open(getConfigPath())
	if err != nil {
		log.Println("Creating new config")
		saveConfig()
		return false
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalln("Unable to read config:", err.Error())
		return false
	}
	splitted := strings.Split(string(data), ";")
	config.host = splitted[0]
	config.password = splitted[1]
	config.serialPort = splitted[2]
	config.dlScriptPath = splitted[3]
	log.Printf("Config loaded; host = %s, serial = %s\n", config.host, config.serialPort)
	return true
}

func resetConfig() bool {
	err := os.Remove(getConfigPath())
	if err != nil {
		log.Println("Unable to reset config?", err.Error())
		return false
	}
	return true
}
