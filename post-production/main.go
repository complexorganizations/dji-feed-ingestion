package main

import (
	"log"
)

// Check if there is a SD card connected.
// Read all the files from the SD card.
// Remove the useless files from the SD card.
// Move all the important data from SD card to local storage.
// Format the SD card; prep it for the next flight.
// Move all the data from the local storage to the youtube.

func main() {
	// Get the mount point of the USB device
	mountPoint := getUSBMountPoint()
	// Get the file path
	filePath := mountPoint + "/"
	// Check if the directory exists
	if !directoryExists(filePath) {
		log.Fatalln("SD card is not connected")
	}
	// Check if the directory is empty
	if isDirectoryEmpty(filePath) {
		log.Fatalln("SD card is empty")
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
	}
	// Name all directories in the directory
	for _, directory := range getAllDirectories {
		// Remove all the empty directories
		if isDirectoryEmpty(directory) {
			removeDirectory(directory)
		}
	}
	// Move all the files from the SD card to the local storage
	newLocation := getCurrentWorkingDirectory() + generateRandomString(10) + "/"
	if !directoryExists(newLocation) {
		createDirectory(newLocation, 0755)
	}
	// Move all the files from the SD card to the local storage, in a new directory with the date and time.
	for _, file := range getAllFiles {
		// Get the file extension
		fileExtension := getFileExtension(file)
		if fileExtension == ".MP4" || fileExtension == ".SRT" {
			moveFile(file, newLocation)
		}
		log.Println("File:", file)
	}
	// Format the SD card
	if directoryExists(filePath) {
		if !isDirectoryEmpty(filePath) {
			nukeDirectory(filePath)
		}
	}
	// Validate if the SD card is empty
	if isDirectoryEmpty(filePath) {
		log.Println("SD card is empty")
	}
}
