package adb

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/alexfacciorusso/adsk/archive"
	"github.com/alexfacciorusso/adsk/envpath"
)

// InstallAdb installs adb when not existing
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
	envpath.AppendToEnvPath(directory + "platform-tools")

	return true
}
