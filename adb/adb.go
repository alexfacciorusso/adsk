package adb

import (
	"adut/archive"
	"adut/register"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/jedib0t/go-pretty/v6/table"
	adb "github.com/zach-klippenstein/goadb"
)

// IsAdbInstalled returns `true` if the adb executable is found in the PATH
func IsAdbInstalled() bool {
	_, err := exec.LookPath("adb")
	return err == nil
}

// InstallAdb ...
func InstallAdb() bool {
	url := "https://dl.google.com/android/repository/platform-tools-latest-windows.zip"
	resp, err := http.Get(url)

	if err != nil {
		log.Printf("Error in getting the zip file %s: %s", url, err)
		return false
	}

	tmpfile, err := ioutil.TempFile("", "platform_tools.zip")

	fmt.Println("Downloading file to ", tmpfile.Name())

	if err != nil {
		log.Printf("Error in creating a temp zip file: %s", err)
		return false
	}

	zipBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error in unzipping the adb zip file: %s", err)
		return false
	}

	_, err = tmpfile.Write(zipBytes)
	if err != nil {
		log.Fatal(err)
	}
	tmpfile.Close()

	directory := "C:/AndroidUtils/"
	archive.Unzip(tmpfile.Name(), directory)
	os.Remove(tmpfile.Name())

	os.Getenv("ANDROID_HOME")
	register.AppendPath(directory + "platform-tools")

	return true
}

// Client wraps the Adb library client
type Client struct {
	adb adb.Adb
}

// NewAdbClient creates a new AdbClient
func NewAdbClient() *Client {
	cli, err := adb.New()

	if err != nil {
		log.Fatal(err)
	}

	err = cli.StartServer()

	if err != nil {
		log.Fatal("Encountered error creating AdbClient: ", err)
	}

	return &Client{adb: *cli}
}

// SetAnimationsEnabled enables or disables the animations on an Android device.
func (a *Client) SetAnimationsEnabled(enabled bool) bool {
	var scale string

	if enabled {
		scale = "1"
	} else {
		scale = "0"
	}

	success := a.execAdbCommand("shell", getAnimationScaleParameters("window", scale)...)
	success = success && a.execAdbCommand("shell", getAnimationScaleParameters("transition", scale)...)
	success = success && a.execAdbCommand("shell", getAnimationScaleParameters("animator", scale)...)

	return success
}

// PrintAllDevices prints all the devices connected to adb in a table.
func (a *Client) PrintAllDevices() {
	devices, err := a.adb.ListDevices()

	if err != nil {
		log.Fatal(err)
	}

	t := table.NewWriter()
	t.SetStyle(table.StyleRounded)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Serial", "Product"})
	t.AppendSeparator()
	for _, device := range devices {
		t.AppendRow(table.Row{device.Serial, device.Product})
	}
	t.Render()
}

func (a *Client) execAdbCommand(cmd string, parameters ...string) bool {
	device := a.adb.Device(adb.AnyDevice())
	_, err := device.RunCommand(cmd, parameters...)

	if err != nil {
		log.Println("Encountered error: ", err)
		return false
	}

	return true
}

func getAnimationScaleParameters(parameter string, scale string) []string {
	return []string{
		"settings", "put", "global", parameter + "_animation_scale", scale,
	}
}
