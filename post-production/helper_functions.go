package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
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

	"github.com/asticode/go-astisub"
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

// Concatenate all subtitles files in a given slice.
func concatenateSubtitlesFiles(fileList []string, outputLocation string, concatenateWaitGroup *sync.WaitGroup) {
	// Sort the files
	sort.Strings(fileList)
	// Create a new empty Subtitles object to hold the combined SRT data
	subtitles := astisub.NewSubtitles()
	// Loop through each SRT file and add its data to the combined Subtitles object
	for _, file := range fileList {
		file, err := astisub.OpenFile(file)
		if err != nil {
			log.Println(err)
		}
		for _, item := range file.Items {
			subtitles.Items = append(subtitles.Items, item)
		}
	}
	// Save the combined SRT data to a new file
	err := subtitles.Write(outputLocation)
	if err != nil {
		log.Println(err)
	}
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

// Get the hmac message authentication
func getHMACMessageAuthentication(content []byte, password []byte) string {
	hash := hmac.New(sha256.New, password)
	hash.Write(content)
	return hex.EncodeToString(hash.Sum(nil))
}

// Validate HMAC for message authentication
func validateHMACMassageAuthentication(content []byte, password []byte, contentSHA string) bool {
	decodedSHA, err := hex.DecodeString(contentSHA)
	if err != nil {
		log.Fatalln(err)
	}
	mac := hmac.New(sha256.New, password)
	mac.Write(content)
	return hmac.Equal(decodedSHA, mac.Sum(nil))
}

/* It takes in a path and content to write to that file.
It uses the os.WriteFile function to write the content to that file.
It checks for errors and logs them. */
func writeToFile(path string, content string) {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}
