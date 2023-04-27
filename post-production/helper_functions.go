package main

import (
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

// Check if there is a usb mounted and if there is return it.
func getUSBMountPoint() string {
	cmd := exec.Command("df", "-h")
	output, err := cmd.Output()
	if err != nil {
		return "Error: " + err.Error()
	}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "/media") || strings.Contains(line, "/run/media") {
			fields := strings.Fields(line)
			return fields[len(fields)-1]
		}
	}
	return "USB device not found"
}

// Walk through a route, find all the files and attach them to a slice.
func walkAndAppendPath(walkPath string) []string {
	var filePath []string
	err := filepath.Walk(walkPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if fileExists(path) {
			filePath = append(filePath, path)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	sort.Strings(filePath)
	return filePath
}

// Walk through a route, find all the folders and attach them to a slice.
func walkAndAppendDirectory(walkPath string) []string {
	var directoryPath []string
	err := filepath.Walk(walkPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if directoryExists(path) {
			directoryPath = append(directoryPath, path)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	sort.Strings(directoryPath)
	return directoryPath
}

// Get the file extension of a file
func getFileExtension(path string) string {
	return filepath.Ext(path)
}

// Check if the file exists and return a bool.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// Checks if the directory exists
func directoryExists(path string) bool {
	directory, err := os.Stat(path)
	if err != nil {
		return false
	}
	return directory.IsDir()
}

// Remove a file from the file system
func removeFile(path string, removeWaitGroup *sync.WaitGroup) {
	err := os.Remove(path)
	if err != nil {
		log.Println(err)
	}
	removeWaitGroup.Done()
}

// Remove a directory and all its contents.
func removeDirectory(dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		log.Println(err)
	}
}

// Check if a given directory is empty.
func isDirectoryEmpty(path string) bool {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Println(err)
	}
	return len(files) == 0
}

// Nuke a given directory and all its contents.
func nukeDirectory(path string) {
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fileExists(path) {
			err := os.Remove(path)
			if err != nil {
				log.Println(err)
			}
		}
		if directoryExists(path) {
			err := os.RemoveAll(path)
			if err != nil {
				log.Println(err)
			}
		}
		return nil
	})
}

// Move a file from one location to another.
func moveFile(source string, destination string, moveWaitGroup *sync.WaitGroup) {
	// Get the file name from the source path
	fileName := filepath.Base(source)
	// Move the file to the destination
	cmd := exec.Command("cp", source, destination+fileName)
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
	moveWaitGroup.Done()
}

// Generate a random string of a given length.
func generateRandomString(length int) string {
	randomBytes := make([]byte, length/2)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("%x", randomBytes)
}

/*
The function takes two parameters: path and permission.
We use os.Mkdir() to create the directory.
If there is an error, we use log.Println() to log the error and then exit the program.
*/
func createDirectory(path string, permission os.FileMode) {
	err := os.Mkdir(path, permission)
	if err != nil {
		log.Println(err)
	}
}

// Get the current working directory on where the executable is running
func getCurrentWorkingDirectory() string {
	currentFileName, err := os.Executable()
	if err != nil {
		log.Println(err)
	}
	return filepath.Dir(currentFileName) + "/"
}

// Get the current time in the format of MM-DD-YYYY_HH-MM-SS and return it as a string.
func getCurrentTime() string {
	return time.Now().Format("01-02-2006_15-04-05")
}

// Lockdown the application to a single linux operating system.
func lockdownToLinuxOperatingSystem() {
	// Check if the operating system is linux
	if runtime.GOOS != "linux" {
		log.Fatalln("This application is only supported on linux operating systems.")
	}
	// Check if the file exists
	validateEtcOsReleaseFileExists := fileExists("/etc/os-release")
	if !validateEtcOsReleaseFileExists {
		log.Fatalln("The file /etc/os-release does not exist.")
	}
	// Read the /etc/os-release file and check if it contains the word "Ubuntu"
	completeEtcOsReleaseFileContent := readAFileAsString("/etc/os-release")
	// Check the name of the operating system
	if strings.Contains(completeEtcOsReleaseFileContent, "ID=ubuntu") {
		// Check the version of the operating system
		if !strings.Contains(completeEtcOsReleaseFileContent, "VERSION_ID=\"22") {
			log.Fatalln("This application is only supported on Ubuntu 22")
		}
	} else {
		log.Fatalln("This application is only supported on Ubuntu.")
	}
}

// Read a file and return the contents
func readAFileAsString(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Println(err)
	}
	return string(content)
}

/*
Imports the "os" package which provides the UserHomeDir() function
Defines the currentUserHomeDir() function
Invokes the UserHomeDir() function
Returns the home directory of the current user
*/
func currentUserHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Println(err)
	}
	if len(homeDir) == 0 {
		homeDir = "/root"
	}
	return homeDir
}

// Get the sha 256 of a file and return it as a string
func sha256OfFile(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
	}
	hash := sha512.New()
	io.Copy(hash, file)
	err = file.Close()
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// Concatenate videos using ffmpeg
func concatenateVideos(videoFiles []string, outputFile string, concatenateWaitGroup *sync.WaitGroup) {
	log.Println("Concatenating videos...")
	log.Println(videoFiles, outputFile)
	inputs := "concat:"
	// Build the input string for ffmpeg command
	for fileIndex, file := range videoFiles {
		// Add the file name to the input string
		inputs = inputs + file
		// Add a pipe if it is not the last file
		if fileIndex < len(videoFiles)-1 {
			inputs = inputs + "|"
		}
	}
	// Execute the ffmpeg command
	cmd := exec.Command("ffmpeg", "-i", inputs, "-c", "copy", outputFile)
	// Run the command
	err := cmd.Run()
	// Check for errors
	if err != nil {
		log.Fatalf("Error running ffmpeg command: %v", err)
	}
	// Remove the files
	for _, file := range videoFiles {
		removeWaitGroup.Add(1)
		go removeFile(file, &removeWaitGroup)
	}
	// Wait for the files to be removed
	removeWaitGroup.Wait()
	// Mark the wait group as done
	concatenateWaitGroup.Done()
}

// Check if the application is installed and in path
func commandExists(application string) bool {
	_, err := exec.LookPath(application)
	return err == nil
}

// Walk through a route, find all the files and attach them to a slice.
func walkAndAppendPathByFileType(walkPath string, fileType string) []string {
	var filePath []string
	err := filepath.Walk(walkPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if fileExists(path) {
			if getFileExtension(path) == fileType {
				filePath = append(filePath, path)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}
	return filePath
}
