package main

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/bluenviron/gortsplib/v3"
	"github.com/bluenviron/gortsplib/v3/pkg/url"
)

/*
It checks if the file exists
If the file exists, it returns true
If the file does not exist, it returns false
*/
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

/*
It takes in a path and content to write to that file.
It uses the os.WriteFile function to write the content to that file.
It checks for errors and logs them.
*/
func writeToFile(path string, content []byte) {
	err := os.WriteFile(path, content, 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

// Read a file and return the contents
func readAFileAsString(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	return string(content)
}

// Get the sha 256 of a file and return it as a string
func sha256OfFile(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalln(err)
	}
	hash := sha512.New()
	io.Copy(hash, file)
	err = file.Close()
	if err != nil {
		log.Fatalln(err)
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// Save all the errors in a single given path.
func saveAllErrors(errors string) {
	// Save the errors in a file
	appendAndWriteToFile(applicationLogFile, errors)
	log.Fatalln(errors)
}

// Check if the application is installed and in path
func commandExists(application string) bool {
	_, err := exec.LookPath(application)
	return err == nil
}

// Check if the json is valid.
func jsonValid(content []byte) bool {
	return json.Valid(content)
}

// Encode struct data to JSON.
func encodeStructToJSON(content interface{}) []byte {
	contentJSON, err := json.Marshal(content)
	if err != nil {
		log.Fatalln(err)
	}
	return contentJSON
}

// Unmarshal json into a struct and return the struct.
func unmarshalJSONIntoStruct(content []byte, data interface{}) interface{} {
	err := json.Unmarshal(content, &data)
	if err != nil {
		log.Fatalln(err)
	}
	return data
}

// Read a file and than return the content as bytes
func readFileAndReturnAsBytes(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
	}
	err = file.Close()
	if err != nil {
		log.Fatalln(err)
	}
	return content
}

// Check if a given rtsp server is alive and responding to requests
func checkRTSPServerAlive(rtspURL string) bool {
	// parse the URL of the server
	parsedURL, err := url.Parse(rtspURL)
	if err != nil {
		return false
	}
	// Connect to the server and close the connection when done
	serverConnection := gortsplib.Client{}
	err = serverConnection.Start(parsedURL.Scheme, parsedURL.Host)
	if err != nil {
		return false
	}
	// Close the connection
	defer serverConnection.Close()
	// Check if the server is alive
	_, _, _, err = serverConnection.Describe(parsedURL)
	return err == nil
}

// Note: Check the packets of the rtsp server in the background.
// Note: If the packets loop than do a counter and end the stream since its a bad stream; recheck and do it again. (loop)

// Run this function in the background and check if a given RTSP server is alive
func checkRTSPServerAliveInBackground(rtspURL string) {
	for {
		// Check if the server is alive
		if checkRTSPServerAlive(rtspURL) {
			rtspServerOneStatus = true
		} else {
			rtspServerOneStatus = false
		}
		// Sleep for 2 seconds, after each check.
		time.Sleep(2 * time.Second)
	}
}

// Forward data to google cloud vertex AI.
func forwardDataToGoogleCloudVertexAI(host string, projectName string, gcpRegion string, vertexStreams string, forwardingWaitGroup *sync.WaitGroup) {
	cmd := exec.Command("vaictl", "-p", projectName, "-l", gcpRegion, "-c", "application-cluster-0", "--service-endpoint", "visionai.googleapis.com", "send", "rtsp", "to", "streams", vertexStreams, "--rtsp-uri", host)
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
	forwardingWaitGroup.Done()
}

// Forward data to AWS Kinesis Video Streams using gstreamer.
func runGstPipeline(host string, streamName string, accessKey string, secretKey string, awsRegion string, forwardingWaitGroup *sync.WaitGroup) {
	cmd := exec.Command("gst-launch-1.0", "rtspsrc", "location="+host, "!", "rtph264depay", "!", "h264parse", "!", "video/x-h264,stream-format=avc", "!", "kvssink", "stream-name="+streamName, "access-key="+accessKey, "secret-key="+secretKey, "aws-region="+awsRegion)
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
	forwardingWaitGroup.Done()
}

// Get the current working directory on where the executable is running
func getCurrentWorkingDirectory() string {
	currentFileName, err := os.Executable()
	if err != nil {
		log.Fatalln(err)
	}
	return filepath.Dir(currentFileName) + "/"
}

// Lockdown the application to a single linux operating system.
func lockdownToLinuxOperatingSystem() {
	// Check if the operating system is linux
	if runtime.GOOS != "linux" {
		saveAllErrors("This application is only supported on linux operating systems.")
	}
	// Check if the file exists
	validateEtcOsReleaseFileExists := fileExists("/etc/os-release")
	if !validateEtcOsReleaseFileExists {
		saveAllErrors("This application is only supported on Ubuntu.")
	}
	// Read the /etc/os-release file and check if it contains the word "Ubuntu"
	completeEtcOsReleaseFileContent := readAFileAsString("/etc/os-release")
	// Check the name of the operating system
	if strings.Contains(completeEtcOsReleaseFileContent, "ID=ubuntu") {
		// Check the version of the operating system
		if !strings.Contains(completeEtcOsReleaseFileContent, "VERSION_ID=\"22.04\"") {
			saveAllErrors("This application is only supported on Ubuntu 22.04.")
		}
	} else {
		saveAllErrors("This application is only supported on Ubuntu.")
	}
}

// Append and write to file
func appendAndWriteToFile(path string, content string) {
	filePath, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = filePath.WriteString(content + "\n")
	if err != nil {
		log.Fatalln(err)
	}
	err = filePath.Close()
	if err != nil {
		log.Fatalln(err)
	}
}

// Get the extension of a given file.
func getFileExtension(path string) string {
	return filepath.Ext(path)
}

// Validate the length of the JSON file.
func validateJSONLength(key string, value string) bool {
	// Check if the key and value are not empty
	if len(key) >= 1 && len(value) >= 1 {
		return true
	}
	// Check if the key and value are not empty
	if len(key) >= 1 && len(value) == 0 {
		saveAllErrors("The value of the key '" + key + "' is empty.")
	}
	return false
}
