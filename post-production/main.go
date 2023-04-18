package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Check if there is a SD card connected.
// Read all the files from the SD card.
// Remove the useless files from the SD card.
// Move all the important data from SD card to local storage.
// Format the SD card; prep it for the next flight.
// Move all the data from the local storage to the youtube.

func main() {
	// Get the mount point of the USB device
	mountPoint, err := getUSBMountPoint()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// Get the file path
	filePath := mountPoint + "/"
	// Check if the directory exists
	if directoryExists(filePath) {
		// Get all files in the directory
		getAllFiles := walkAndAppendPath(filePath)
		// Name all files in the directory
		for _, file := range getAllFiles {
			// Get the file extension
			fileExtension := getFileExtension(file)
			// Remove all files that are not MP4 or SRT
			if fileExtension != ".MP4" && fileExtension != ".SRT" {
				removeFile(file)
			}
			// Get the file path
			log.Println("File:", file)
		}
	}
}

func getUSBMountPoint() (string, error) {
	cmd := exec.Command("df", "-h")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "/media") || strings.Contains(line, "/run/media") {
			fields := strings.Fields(line)
			return fields[len(fields)-1], nil
		}
	}
	return "", fmt.Errorf("USB device not found")
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
		log.Fatalln(err)
	}
	return filePath
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
func removeFile(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Fatalln(err)
	}
}

// Remove a directory and all its contents.
func removeDirectory(dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		log.Fatalln(err)
	}
}

/* It takes the path of a directory as an argument.
If the directory is empty, it returns a true value.
Otherwise, it returns a false value. */
func isDirectoryEmpty(path string) bool {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatalln(err)
	}
	return len(files) == 0
}
