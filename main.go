package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

func main() {
	reloaderRunning := true

	if len(os.Args) < 2 {
		usage()
		return
	}

	path := os.Args[1]

	interruptChan := make(chan os.Signal)
	signal.Notify(interruptChan, os.Interrupt)

	for reloaderRunning {

		runChan := make(chan bool)
		cmd := run(runChan)

		modTimeChan := make(chan bool)
		go watchModTime(path, modTimeChan)

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
					log.Print("Done!")
				}
			case msg := <-modTimeChan:
				if !msg {
					cmd.Process.Signal(os.Interrupt)
					subprocessRunning = false
					reloaderRunning = false
				} else {
					log.Print("Executable modified. Restarting.")
					cmd.Process.Signal(os.Interrupt)
					restarting = true
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
func run(runChan chan bool) *exec.Cmd {
	cmd := exec.Command(os.Args[1], os.Args[2:len(os.Args)]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		log.Fatal("Process start failed.")
	}

	go func() {
		cmd.Wait()
		runChan <- true
	}()

	return cmd
}

// watches the given path for a change based on the mtime. Signals true
// on modChan if a change happens, or false if the file can't be checked
// for mtime for some reason
func watchModTime(path string, modChan chan bool) {
	initialModTime, error := modTime(path)
	if error != nil {
		log.Fatal("Could not determine mtime for path")
		modChan <- false
		return
	}

	for {
		newModTime, error := modTime(path)
		if error != nil {
			log.Fatal("Could not determine mtime for path")
			modChan <- false
			return
		}

		if newModTime.After(initialModTime) {
			modChan <- true
			return
		}

		time.Sleep(1000)
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
