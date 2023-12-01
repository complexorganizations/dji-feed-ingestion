package main

import (
	"bufio"
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
	"strconv"
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
		if strings.Contains(line, "Mavic-3-") {
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

// Get the current time in the format of MM-DD-YYYY and return it as a string.
func getCurrentTime() string {
	return time.Now().Format("01-02-2006")
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

// Compress a video video file using ffmpeg
func compressVideo(videoFile string, compressWaitGroup *sync.WaitGroup) {
	log.Println("Compressing video file: " + videoFile)
	// Get the file path only from the video file
	fileDirectoryOnly := videoFile[0:strings.LastIndex(videoFile, "/")]
	// Get the file name only from the file path, remove the extension
	fileNameOnly := videoFile[strings.LastIndex(videoFile, "/")+1 : strings.LastIndex(videoFile, ".")]
	// Create the compressed video file name
	compressedVideoFile := fileDirectoryOnly + "/" + fileNameOnly + "_compressed.mp4"
	// Compress the video file
	cmd := exec.Command("ffmpeg", "-i", videoFile, "-vcodec", "libx265", "-crf", "23", compressedVideoFile)
	// Run the command
	err := cmd.Run()
	// Check if there was an error
	if err != nil {
		log.Println(err)
	}
	// Remove the original video file
	os.Remove(videoFile)
	// Close the wait group
	compressWaitGroup.Done()
}

// Concatenate videos using ffmpeg
func concatenateVideos(videoFiles []string, outputFile string, concatenateWaitGroup *sync.WaitGroup) {
	/*
		// Compress all the video files so they are smaller; combine the video files into one.
		compressWaitGroup := sync.WaitGroup{}
		for _, videoFile := range videoFiles {
			compressWaitGroup.Add(1)
			go compressVideo(videoFile, &compressWaitGroup)
		}
		// Wait for all the video files to be compressed
		compressWaitGroup.Wait()
	*/
	// Concatenate the video files
	log.Println("Concatenating videos...")
	log.Println(videoFiles, outputFile)
	// Get the output directory
	outputDirectory := filepath.Dir(outputFile)
	// NOTE: Compress all the video files so they are smaller; combine the video files into one.
	// Create a temp path to store the video file log
	tempVideoFilesPath := outputDirectory + "/tempCombineVideos.txt"
	// Write the input string for ffmpeg command
	videoConcatenateWriteFile(videoFiles, tempVideoFilesPath)
	// Execute the ffmpeg command
	cmd := exec.Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", tempVideoFilesPath, "-c", "copy", outputFile)
	// Run the command
	err := cmd.Run()
	// Check for errors
	if err != nil {
		log.Fatalf("Error running ffmpeg command: %v", err)
	}
	// Remove the temp file
	removeWaitGroup.Add(1)
	go removeFile(tempVideoFilesPath, &removeWaitGroup)
	// Remove the files
	for _, file := range videoFiles {
		removeWaitGroup.Add(1)
		go removeFile(file, &removeWaitGroup)
	}
	// Wait for the files to be removed
	removeWaitGroup.Wait()
	log.Println("Concatenation complete.")
	// Mark the wait group as done
	concatenateWaitGroup.Done()
}

// Write the list of video files to concatenate using ffmpeg
func videoConcatenateWriteFile(videoFiles []string, fileName string) {
	// Open the file for writing
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	// Write each video file name to the file in the required format
	for _, videoFile := range videoFiles {
		_, err := file.WriteString(fmt.Sprintf("file '%s'\n", videoFile))
		if err != nil {
			log.Fatalln(err)
		}
	}
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

// SRTEntry defines the structure of a subtitle entry in an SRT file.
type SRTEntry struct {
	Index     int    // The sequence number of the subtitle
	StartTime string // The start time of the subtitle display
	EndTime   string // The end time of the subtitle display
	Text      string // The actual subtitle text
}

// Concatenate all subtitles files in a given slice.
func concatenateSubtitlesFiles(fileList []string, outputLocation string, concatenateWaitGroup *sync.WaitGroup) {
	// Sort the files
	sort.Strings(fileList)
	// Check if there are enough filenames
	if len(fileList) < 2 {
		log.Println("Please provide at least two SRT files.")
	}
	// Read the first SRT file
	srts := readSRT(fileList[0])
	// Loop over all the remaining filenames
	for _, filename := range fileList[1:] {
		// Read the current SRT file
		currentSRT := readSRT(filename)
		// Combine the current SRT file with the combined SRT files
		srts = combineSRT(srts, currentSRT)
	}
	// Write the combined SRT data into a new file
	writeSRT(outputLocation, srts)
	// Remove the files
	for _, file := range fileList {
		removeWaitGroup.Add(1)
		go removeFile(file, &removeWaitGroup)
	}
	// Wait for the files to be removed
	removeWaitGroup.Wait()
	// Mark the wait group as done
	concatenateWaitGroup.Done()
}

// Function to read an SRT file and return a slice of SRT entries
func readSRT(filename string) []SRTEntry {
	// Open the file
	file, err := os.Open(filename)
	// Print the error if unable to open the file
	if err != nil {
		log.Println(err)
	}
	// Create a new scanner for the file
	scanner := bufio.NewScanner(file)
	// Set the split function for the scanning operation
	scanner.Split(splitSRTEntries)
	// Create an empty slice to hold SRT entries
	var entries []SRTEntry
	// Loop over all the lines in the file
	for scanner.Scan() {
		// Parse the current line into an SRT entry
		entry := parseSRTEntry(scanner.Text())
		// Append the SRT entry to the slice of entries
		entries = append(entries, entry)
	}
	// Close the file
	err = file.Close()
	if err != nil {
		// Print the error if unable to close the file
		log.Println(err)
	}
	// Return the slice of SRT entries
	return entries
}

// Function to parse a single line of an SRT file into an SRTEntry struct
func parseSRTEntry(line string) SRTEntry {
	// Split the line into parts by newline
	parts := strings.Split(line, "\n")
	// If there are less than 4 parts, return an empty SRTEntry
	if len(parts) < 4 {
		return SRTEntry{}
	}
	// Convert the first part to an integer for the index
	index, err := strconv.Atoi(parts[0])
	if err != nil {
		// If unable to convert, return an empty SRTEntry
		return SRTEntry{}
	}
	// Split the second part into start and end times
	timeParts := strings.Split(parts[1], " --> ")
	// If there are not 2 time parts, return an empty SRTEntry
	if len(timeParts) != 2 {
		return SRTEntry{}
	}
	// Return an SRTEntry with the parsed data
	return SRTEntry{
		Index:     index,
		StartTime: timeParts[0],
		EndTime:   timeParts[1],
		Text:      strings.Join(parts[2:], "\n"),
	}
}

// Define a new function to split the SRT entries
func splitSRTEntries(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// If at the end of the file and no data left, return nothing
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	// If found a double newline (which separates SRT entries), return the entry
	if i := strings.Index(string(data), "\n\n"); i >= 0 {
		return i + 2, data[0:i], nil
	}
	// If at the end of the file, return the remaining data
	if atEOF {
		return len(data), data, nil
	}
	// Otherwise, return nothing
	return 0, nil, nil
}

// Define a new function to combine two SRT files
func combineSRT(srt1, srt2 []SRTEntry) []SRTEntry {
	// Get the last index of the first SRT file as the offset
	offset := srt1[len(srt1)-1].Index
	// Create an empty slice for the combined SRT entries
	var combined []SRTEntry
	// Parse the end time of the last entry in the first SRT file
	lastEndTime1, _ := time.Parse("15:04:05.000", srt1[len(srt1)-1].EndTime)
	// Parse the start time of the first entry in the second SRT file
	firstStartTime2, _ := time.Parse("15:04:05.000", srt2[0].StartTime)
	// Calculate the time difference between the two
	diffTime := lastEndTime1.Sub(firstStartTime2)
	// Set the initial frame count offset to be the last index of the first SRT file
	frameCountOffset := srt1[len(srt1)-1].Index
	// Loop over all entries in the second SRT file
	for _, entry := range srt2 {
		// Adjust the index of each entry by the offset
		entry.Index = entry.Index + offset
		// Parse the start and end times of the current entry
		entryStartTime, _ := time.Parse("15:04:05.000", entry.StartTime)
		entryEndTime, _ := time.Parse("15:04:05.000", entry.EndTime)
		// Adjust the start and end times by the time difference
		newStartTime := entryStartTime.Add(diffTime)
		newEndTime := entryEndTime.Add(diffTime)
		// Format the new start and end times
		entry.StartTime = newStartTime.Format("15:04:05.000")
		entry.EndTime = newEndTime.Format("15:04:05.000")
		// Increase the frame count offset
		frameCountOffset = frameCountOffset + 1
		// Replace the old frame count with the new frame count in the entry text
		entry.Text = strings.Replace(entry.Text, fmt.Sprintf("FrameCnt: %d", entry.Index-offset), fmt.Sprintf("FrameCnt: %d", frameCountOffset), -1)
		// Append the adjusted entry to the combined slice
		combined = append(combined, entry)
	}
	// Append all entries from the first SRT file to the combined slice
	combined = append(srt1, combined...)
	// Return the combined slice
	return combined

}

// Function to write the combined SRT entries into a new file
func writeSRT(filename string, entries []SRTEntry) {
	// Create a new file with the specified filename
	file, err := os.Create(filename)
	if err != nil {
		// Log the error if the file cannot be created
		log.Println(err)
	}
	// Create a new buffered writer for the file
	writer := bufio.NewWriter(file)
	// Loop over all entries
	for _, entry := range entries {
		// Write each entry as a formatted string to the file
		_, err := writer.WriteString(fmt.Sprintf("%d\n%s --> %s\n%s\n\n", entry.Index, entry.StartTime, entry.EndTime, entry.Text))
		if err != nil {
			// Log the error if the entry cannot be written to the file
			log.Println(err)
		}
	}
	// Flush any remaining data in the buffer to the file
	err = writer.Flush()
	if err != nil {
		// Log the error if the data cannot be flushed to the file
		log.Println(err)
	}
	// Close the file
	err = file.Close()
	if err != nil {
		// Log the error if the file cannot be closed
		log.Println(err)
	}
}
