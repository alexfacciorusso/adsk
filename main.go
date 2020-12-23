package main

import (
	"adut/adb"
	"adut/ostools"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/logrusorgru/aurora"
	"github.com/manifoldco/promptui"
)

type logPrinter struct{}

func (l *logPrinter) Write(p []byte) (int, error) {
	return color.Error.Write([]byte(color.YellowString(string(p))))
}

func main() {
	var verboseFlag bool
	flag.BoolVar(&verboseFlag, "v", false, "If set, debug messages will be print.")

	flag.Parse()

	if verboseFlag {
		log.SetOutput(new(logPrinter))
	} else {
		log.SetOutput(ioutil.Discard)
	}

	isAdbInstalled := adb.IsAdbInstalled()
	if !isAdbInstalled {
		prompt := promptui.Prompt{
			Label:     "I can't find 'adb' in your system's PATH, but I can install it for you. Do you want me to do so",
			IsConfirm: true,
		}

		answer, err := prompt.Run()

		if err != nil {
			fmt.Printf("\nPrompt failed %v\n", err)
			return
		}

		if strings.ToLower(answer) != "y" {
			fmt.Println("\nI'm sorry, but without adb I can't work properly. I will need to die now. ", aurora.Magenta("Thank you."))
			os.Exit(1)
		}

		if !ostools.AmAdmin() {
			fmt.Println("\nI am not running as admin, I need to restart in order to install adb. Please give me access to admin when requested")
			ostools.RunMeElevated()
		}
		adb.InstallAdb()
	}

	adbClient := adb.NewAdbClient()

	exit := false

	for !exit {
		choice := promptui.Select{
			Label: "What action do you want me to do?",
			Items: []string{
				"List all devices connected to adb",
				"Disable your device's animations for testing",
				"Enable your device's animations for testing",
				"Exit",
			},
			Templates: &promptui.SelectTemplates{
				Label:    `{{ "[?]" | blue }} {{ . }}`,
				Active:   "> {{ . | yellow }}",
				Selected: "{{ . | yellow | bold}}",
			},
		}

		i, _, err := choice.Run()

		if err != nil {
			os.Exit(1)
		}

		switch i {
		case 0:
			adbClient.PrintAllDevices()
		case 1:
			setAnimations(adbClient, false)
		case 2:
			setAnimations(adbClient, true)
		default:
			exit = true
		}
	}
}

func setAnimations(client *adb.Client, enabled bool) bool {
	success := client.SetAnimationsEnabled(enabled)

	if !success {
		color.HiRed("\nI couldn't make it, sorry. Either your device is not connected or there was some problem.\n\n")
	}

	return success
}
