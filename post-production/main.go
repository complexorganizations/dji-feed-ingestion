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
// Move all the data from the local storage to the s3.

func main() {
	mountPoint, err := getUSBMountPoint()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	filePath := mountPoint + "/"
	// Get all files in the directory
	getAllFiles := walkAndAppendPath(filePath)
	// Name all files in the directory
	for _, file := range getAllFiles {
		// Get the file extension
		fileExtension := getFileExtension(file)
		log.Println("extension:", fileExtension)
		if fileExtension != ".MP4" && fileExtension != ".SRT" {
			log.Println("Removing file:", file)
			err := os.Remove(file)
			if err != nil {
				log.Println("Error:", err)
			}
		}
		// Get the file name
		fileName := filepath.Base(file)
		log.Println("name:", fileName)
		// Get the file name without the extension
		fileNameWithoutExtension := strings.TrimSuffix(fileName, fileExtension)
		log.Println("name without extension:", fileNameWithoutExtension)
		// Get the file path
		filePath := filepath.Dir(file)
		log.Println("path:", filePath)
		log.Println("File:", file)
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

// Check if the file exists and return a bool.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// Get the file extension of a file
func getFileExtension(path string) string {
	return filepath.Ext(path)
}
