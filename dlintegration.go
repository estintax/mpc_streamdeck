package main

import (
	"log"
	"os"
	"time"

	"github.com/andreykaipov/goobs/api/requests/inputs"
	"github.com/andreykaipov/goobs/api/requests/scenes"
	"github.com/andreykaipov/goobs/api/requests/ui"
	"github.com/estintax/mpc_streamdeck/dinolang"
)

func initDL() {
	dinolang.Classes["deck"] = &dinolang.Class{
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
				info, err := client.Ui.GetStudioModeEnabled()
				if err == nil && info.StudioModeEnabled {
					params := scenes.NewSetCurrentPreviewSceneParams()
					params.SceneName = &sceneName
					_, err = client.Scenes.SetCurrentPreviewScene(params)
				} else {
					params := scenes.NewSetCurrentProgramSceneParams()
					params.SceneName = &sceneName
					_, err = client.Scenes.SetCurrentProgramScene(params)
				}
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
	case "switch-mute":
		if len(args) > 1 {
			if dinolang.GetTypeEx(args[1]) == "string" {
				audioName := dinolang.StringToText(dinolang.IfVariableReplaceIt(args[1]).(string))
				params := inputs.NewGetInputMuteParams()
				params.InputName = &audioName
				res, err := client.Inputs.GetInputMute(params)
				if err != nil {
					dinolang.PrintError(err.Error())
					dinolang.SetReturned("int", 0, segmentName)
					return false
				}
				muted := !res.InputMuted
				newParams := inputs.NewSetInputMuteParams()
				newParams.InputName = &audioName
				newParams.InputMuted = &muted
				_, err = client.Inputs.SetInputMute(newParams)
				if err != nil {
					dinolang.PrintError(err.Error())
					dinolang.SetReturned("int", 0, segmentName)
					return false
				}
				dinolang.SetReturned("int", 1, segmentName)
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
	case "inputs":
		if len(args) > 1 {
			if dinolang.GetTypeEx(args[1]) == "unknown" {
				res, err := client.Inputs.GetInputList()
				if err != nil {
					dinolang.PrintError(err.Error())
					dinolang.SetReturned("int", 0, segmentName)
					return false
				}
				var inputs []string
				for _, i := range res.Inputs {
					inputs = append(inputs, i.InputName)
				}
				dinolang.SetVariable(args[1], inputs)
				dinolang.SetReturned("int", len(inputs), segmentName)
			} else {
				dinolang.PrintError("Type of the first argument must be unknown")
				dinolang.SetReturned("int", 0, segmentName)
				return false
			}
		} else {
			dinolang.PrintError("Too few arguments")
			dinolang.SetReturned("int", 0, segmentName)
			return false
		}
	case "switch-studio":
		info, err := client.Ui.GetStudioModeEnabled()
		if err != nil {
			dinolang.PrintError(err.Error())
			dinolang.SetReturned("int", 0, segmentName)
			return false
		}
		enabled := true
		if info.StudioModeEnabled {
			enabled = false
			previewInfo, _ := client.Scenes.GetCurrentPreviewScene()
			currentInfo, _ := client.Scenes.GetCurrentProgramScene()
			if previewInfo.CurrentPreviewSceneName != currentInfo.CurrentProgramSceneName {
				res, _ := client.Transitions.GetCurrentSceneTransition()
				client.Scenes.SetCurrentProgramScene(&scenes.SetCurrentProgramSceneParams{SceneName: &previewInfo.CurrentPreviewSceneName})
				time.Sleep(time.Duration(time.Duration(res.TransitionDuration)*time.Millisecond) + 500*time.Millisecond)
			}
		}
		_, err = client.Ui.SetStudioModeEnabled(&ui.SetStudioModeEnabledParams{StudioModeEnabled: &enabled})
		if err != nil {
			dinolang.PrintError(err.Error())
			dinolang.SetReturned("int", 0, segmentName)
			return false
		}
	default:
		dinolang.PrintError("Unknown method")
		return false
	}

	return true
}
