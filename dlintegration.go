package main

import (
	"log"
	"os"

	"github.com/andreykaipov/goobs/api/requests/scenes"
	"github.com/estintax/mpc_streamdeck/dinolang"
)

func initDL() {
	dinolang.Classes["deck"] = dinolang.Class{
		Prefix:    "deck",
		Used:      true,
		IsBuiltIn: false,
		Caller:    DeckClassHandler,
		Loader:    nil}
	if len(os.Args) > 1 && os.Args[1] == "--dl-cli" {
		dinolang.PiniginShell()
	}
}

func runScript(key int) bool {
	dinolang.SetVariable("key", key)
	if !dinolang.ParseFile(config.dlScriptPath) {
		log.Println("Unable to run Dinolang script!")
		reinitAnswer()
		return false
	}
	return true
}

func DeckClassHandler(args []string, segmentName string) bool {
	switch args[0] {
	case "switch-scene":
		if len(args) > 1 {
			if dinolang.GetTypeEx(args[1]) == "string" {
				sceneName := dinolang.StringToText(dinolang.IfVariableReplaceIt(args[1]).(string))
				params := scenes.NewSetCurrentProgramSceneParams()
				params.SceneName = &sceneName
				_, err := client.Scenes.SetCurrentProgramScene(params)
				if err != nil {
					dinolang.PrintError(err.Error())
					dinolang.SetReturned("int", 0, segmentName)
					return false
				}
				dinolang.SetReturned("int", 1, segmentName)
				return false
			} else {
				dinolang.PrintError("Type of the first argument is not a string")
				dinolang.SetReturned("int", 0, segmentName)
				return false
			}
		} else {
			dinolang.PrintError("Too few arguments")
			dinolang.SetReturned("int", 0, segmentName)
			return false
		}
	}

	return true
}
