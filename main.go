package main

import (
	//	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

func main() {
	reloaderRunning := true

	// if not enough arguments were given, go ahead and show the usage and
	// exit
	if len(os.Args) < 2 {
		usage()
		return
	}

	path := os.Args[1]

	// run a quick initial check just to ensure that the binary path is valid
	interruptChan := make(chan os.Signal)
	signal.Notify(interruptChan, os.Interrupt)

	// start the mtime watcher
	modTimeChan := make(chan bool)
	go watchModTime(path, modTimeChan)

	for reloaderRunning {
		runChan := make(chan bool)
		cmd, error := run(runChan)
		if error != nil {
			log.Fatal(error)
		}

		subprocessRunning := true
		restarting := false
		for subprocessRunning {
			select {
			case <-interruptChan:
				// pass interrupt through to the subprocess
				cmd.Process.Signal(os.Interrupt)
			case <-runChan:
				if restarting {
					subprocessRunning = false
				} else {
					subprocessRunning = false
					reloaderRunning = false
					log.Print("Executable exited. Exiting.")
				}
			case msg := <-modTimeChan:
				if !msg {
					log.Print("Error loading mtime for executable.")
				} else {
					cmd.Process.Signal(os.Interrupt)
					restarting = true
					log.Print("Executable modified. Restarting.")
				}
			}
		}
	}
}

// Displays a helpful usage message
func usage() {
	fmt.Println("reloader - A simple executable reloader")
	fmt.Println("Usage: reloader [executable path] [args]")
	fmt.Println()
}

// Runs the command line args after the first one (this command),
func run(runChan chan bool) (*exec.Cmd, error) {
	cmd := exec.Command(os.Args[1], os.Args[2:len(os.Args)]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	retries := 5
	for retry := 1; retry <= retries; retry++ {
		err := cmd.Start()
		if err != nil {
			if retry == retries {
				return nil, err
			}
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}

	go func() {
		cmd.Wait()
		runChan <- true
	}()

	return cmd, nil
}

// watches the given path for a change based on the mtime. Signals true
// on modChan if a change happens, or false if the file can't be checked
// for mtime for some reason
func watchModTime(path string, modChan chan bool) {
	currentModTime, error := modTime(path)
	if error != nil {
		modChan <- false
		return
	}

	for {
		newModTime, error := modTime(path)
		if error != nil {
			modChan <- false
		}

		if !newModTime.Equal(currentModTime) {
			currentModTime = newModTime
			modChan <- true
		}

		time.Sleep(1 * time.Second)
	}
}

// finds the last modified time for the given path, or
// an error if we can't determine that time for some reason
func modTime(path string) (time.Time, error) {
	fileInfo, error := os.Stat(path)
	if error != nil {
		return time.Time{}, error
	}

	return fileInfo.ModTime(), error
}
