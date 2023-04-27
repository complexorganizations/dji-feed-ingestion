package main

import (
	"log"
	"sync"
	"time"
)

var removeWaitGroup sync.WaitGroup
var moveWaitGroup sync.WaitGroup
var concatenateWaitGroup sync.WaitGroup

func init() {
	// Lockdown to the linux OS.
	lockdownToLinuxOperatingSystem()
	// Check if the required programs are installed.
	//
}

// Check if there is a SD card connected.
// Read all the files from the SD card.
// Remove the useless files from the SD card.
// Move all the important data from SD card to local storage.
// Format the SD card; prep it for the next flight.
// Move all the data from the local storage to the youtube.

func main() {
	for {
		// Get the mount point of the USB device
		mountPoint := getUSBMountPoint()
		// Get the file path
		filePath := mountPoint + "/"
		// Check if the SD card is connected
		if directoryExists(filePath) {
			// Check if the SD card is empty
			// Note: Add checks here to check specfic usb and not all drives for this.
			if !isDirectoryEmpty(filePath) {
				// Get all the MP4 files in the directory
				var sdCardVideoFilesOnly []string = walkAndAppendPathByFileType(filePath, ".MP4")
				// Get all the SRT files in the directory
				var sdCardSRTFilesOnly []string = walkAndAppendPathByFileType(filePath, ".SRT")
				// Const file names
				randomFileName := generateRandomString(10) + "_" + getCurrentTime()
				// Move all the files from the SD card to the local storage
				newLocation := getCurrentWorkingDirectory() + randomFileName + "/"
				if !directoryExists(newLocation) {
					createDirectory(newLocation, 0755)
				}
				// Move all the files from the SD card to the local storage, in a new directory with the date and time.
				for _, file := range sdCardVideoFilesOnly {
					moveWaitGroup.Add(1)
					go moveFile(file, newLocation, &moveWaitGroup)
					log.Println("Moving file: " + file)
				}
				for _, file := range sdCardSRTFilesOnly {
					moveWaitGroup.Add(1)
					go moveFile(file, newLocation, &moveWaitGroup)
					log.Println("Moving file: " + file)
				}
				// Wait for all the files to be moved
				moveWaitGroup.Wait()
				// Format the SD card
				if !isDirectoryEmpty(filePath) {
					nukeDirectory(filePath)
				}
				// Start the post processing on the local system here, as a go routine so that it can continue with the loop.
				var videoFilesOnly []string = walkAndAppendPathByFileType(newLocation, ".MP4")
				// Create the variable to store the srt files
				var srtFilesOnly []string = walkAndAppendPathByFileType(newLocation, ".SRT")
				// Create a location to store the final video
				finalVideoLocation := newLocation + randomFileName + ".mp4"
				// Create a location to store the final srt file
				finalSRTLocation := newLocation + randomFileName + ".srt"
				// Add one to the wait group
				concatenateWaitGroup.Add(2)
				// Concatenate all the videos
				go concatenateVideos(videoFilesOnly, finalVideoLocation, &concatenateWaitGroup)
				// Concatenate all the srt files
				go concatenateSubtitlesFiles(srtFilesOnly, finalSRTLocation, &concatenateWaitGroup)
			} else {
				log.Println("SD card is empty.")
			}
		} else {
			log.Println("SD card not found.")
		}
		// Wait 5 seconds before checking again
		time.Sleep(5 * time.Second)
		// Wait for all the files to be concatenated
		concatenateWaitGroup.Wait()
	}
}
