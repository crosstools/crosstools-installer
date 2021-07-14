package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const (
	CURRENT_VERSION_FILE = "crosstools-installer-version.data"
	DIRECTORY_NAME       = "crosstools-installer"
)

var EXECUTABLE_NAME string = "crosstools-installer"

func programFolderToSaveTo() string {
	switch runtime.GOOS {
	case "windows":
		return "C:\\Program Files\\"
	default:
		log.Fatalln("The Operating System you are on is not supported by CrossTools Installer")
	}

	return ""
}

func setPath(path string) {
	switch runtime.GOOS {
	case "windows":
		output, err := exec.Command("Path").Output()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to get PATH info to check for existing crosstools-installer")
			log.Fatalln(err)
		}
		outputStr := strings.Replace(string(output), "PATH=", "", 1) // remove "PATH=" out of output
		paths := strings.Split(outputStr, ";")
		for _, pathFromPaths := range paths {
			if path == pathFromPaths {
				return // path has already been set, so end this function
			}
		}
		err = exec.Command("setx Path \"%Path%;" + path + "\" /m").Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to set PATH for crosstools-installer")
			log.Fatalln(err)
		}
	default:
		log.Fatalln("Path for crosstools-installer has not been set")
	}
}

func InstallSelf() {
	check := func(e error) {
		if e != nil {
			fmt.Fprintln(os.Stderr, "Install self error: Failed to install crosstools-installer self")
			panic(e)
		}
	}

	if runtime.GOOS == "windows" {
		EXECUTABLE_NAME += ".exe"
	}

	directory := programFolderToSaveTo() + DIRECTORY_NAME
	resp, err := http.Get("https://raw.githubusercontent.com/crosstools/crosstools-installer/main/crosstools-installer.exe")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get executable of crosstools-installer")
		check(err)
	}
	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read executable response of crosstools-installer")
		check(err)
	}

	respVersion, err := http.Get("https://raw.githubusercontent.com/crosstools/crosstools-installer/main/VERSION")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get version of crosstools-installer")
		check(err)
	}
	defer respVersion.Body.Close()

	respVersionData, err := ioutil.ReadAll(respVersion.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read data of version for crosstools-installer")
		check(err)
	}

	check(os.RemoveAll(directory))
	check(os.Mkdir(directory, 0755))

	err = ioutil.WriteFile(directory+string(os.PathSeparator)+CURRENT_VERSION_FILE, respVersionData, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to write to '"+CURRENT_VERSION_FILE+"' file")
		check(err)
	}

	err = ioutil.WriteFile(directory+string(os.PathSeparator)+EXECUTABLE_NAME, respData, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to write to '"+EXECUTABLE_NAME+"' file")
		check(err)
	}

	// Set PATH
	setPath(directory + string(os.PathSeparator) + EXECUTABLE_NAME)

	fmt.Println("Successfully installed crosstools-installer, enjoy using it!")
}

// func UpdateSelf() {
// 	resp, err := http.Get("https://raw.githubusercontent.com/crosstools/crosstools-installer/main/VERSION")
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, "Updating self error:", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	respData, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, "Updating self reading response error:", err)
// 		return
// 	}

// 	version := string(respData)

// 	if _, err := os.Stat(CURRENT_VERSION_FILE); os.IsNotExist(err) {
// 		fmt.Fprintf(os.Stderr, "File '%s' does not exist so we don't know what version crosstools-installer is in.\nWe will install the latest version anyways, next time please do not delete the '%s' file.\n", CURRENT_VERSION_FILE, CURRENT_VERSION_FILE)
// 		os.PathSeparator
// 	}
// }
