package main

import (
	"fmt"
	"log"
	"os"
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
	if !directoryExists(filePath) {
		os.Exit(1)
	}
	// Get all files in the directory
	getAllFiles := walkAndAppendPath(filePath)
	// Get all directories in the directory
	getAllDirectories := walkAndAppendDirectory(filePath)
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
	// Name all directories in the directory
	for _, directory := range getAllDirectories {
		// Remove all the empty directories
		if isDirectoryEmpty(directory) {
			removeDirectory(directory)
		}
		// Get the directory path
		log.Println("Directory:", directory)
	}

	// Move all the files from the SD card to the local storage, in a new directory with the date and time.
	for _, file := range getAllFiles {
		// Move all the files from the SD card to the local storage
		moveFile(file, "/home/prajwal/Projects/dji-feed-analysis/post-production/")
	}
	// Remove the directory
	removeAllFilesInDirectory(filePath)
	if isDirectoryEmpty(filePath) {
		log.Println("Directory is empty")
	}
}
