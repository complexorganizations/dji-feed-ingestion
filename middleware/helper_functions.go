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

var (
	mutex = &sync.RWMutex{}
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
		log.Println(err)
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
		log.Println(err)
	}
	return contentJSON
}

// Unmarshal json into a struct and return the struct.
func unmarshalJSONIntoStruct(content []byte, data interface{}) interface{} {
	err := json.Unmarshal(content, &data)
	if err != nil {
		log.Println(err)
	}
	return data
}

// Read a file and than return the content as bytes
func readFileAndReturnAsBytes(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}
	content, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
	}
	err = file.Close()
	if err != nil {
		log.Println(err)
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
			mutex.Lock()
			if getValueFromMap(rtspServerStatusChannel, rtspURL) == false {
				addKeyValueToMap(rtspServerStatusChannel, rtspURL, true)
			}
			mutex.Unlock()
		} else {
			mutex.Lock()
			if getValueFromMap(rtspServerStatusChannel, rtspURL) == true {
				addKeyValueToMap(rtspServerStatusChannel, rtspURL, false)
			}
			mutex.Unlock()
		}
		// Sleep for 3 seconds, after each check.
		time.Sleep(3 * time.Second)
	}
}

// Forward data to google cloud vertex AI.
func forwardDataToGoogleCloudVertexAI(host string, projectName string, gcpRegion string, vertexStreams string, forwardingWaitGroup *sync.WaitGroup) {
	// Set the rtspServerStreamingChannel to true
	rtspServerStreamingChannel[host] = true
	// Move the default file to a temporary file.
	if fileExists(amazonKinesisDefaultPath) {
		moveFile(amazonKinesisDefaultPath, amazonKinesisTempPath)
	}
	// Run the command to forward the data to vertex AI
	cmd := exec.Command("vaictl", "-p", projectName, "-l", gcpRegion, "-c", "application-cluster-0", "--service-endpoint", "visionai.googleapis.com", "send", "rtsp", "to", "streams", vertexStreams, "--rtsp-uri", host)
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
	// Once the data is forwarded, remove the temporary file.
	if fileExists(amazonKinesisTempPath) {
		moveFile(amazonKinesisTempPath, amazonKinesisDefaultPath)
	}
	// Done forwarding
	forwardingWaitGroup.Done()
	// Set the rtspServerStreamingChannel to false
	rtspServerStreamingChannel[host] = false
}

// Forward data to AWS Kinesis Video Streams using gstreamer.
func forwardDataToAmazonKinesisStreams(host string, streamName string, accessKey string, secretKey string, awsRegion string, forwardingWaitGroup *sync.WaitGroup) {
	// Set the rtspServerStreamingChannel to true
	rtspServerStreamingChannel[host] = true
	// Move the temporary file to the default file location if it exists.
	if fileExists(amazonKinesisTempPath) {
		moveFile(amazonKinesisTempPath, amazonKinesisDefaultPath)
	}
	// Run the gstreamer command to forward the data to AWS Kinesis Video Streams
	cmd := exec.Command("gst-launch-1.0", "rtspsrc", "location="+host, "!", "rtph264depay", "!", "h264parse", "!", "video/x-h264,stream-format=avc", "!", "kvssink", "stream-name="+streamName, "access-key="+accessKey, "secret-key="+secretKey, "aws-region="+awsRegion)
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
	forwardingWaitGroup.Done()
	// Set the rtspServerStreamingChannel to false
	rtspServerStreamingChannel[host] = false
}

// Get the current working directory on where the executable is running
func getCurrentWorkingDirectory() string {
	currentFileName, err := os.Executable()
	if err != nil {
		log.Println(err)
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
		saveAllErrors("The file /etc/os-release does not exist.")
	}
	// Read the /etc/os-release file and check if it contains the word "Ubuntu"
	completeEtcOsReleaseFileContent := readAFileAsString("/etc/os-release")
	// Check the name of the operating system
	if strings.Contains(completeEtcOsReleaseFileContent, "ID=ubuntu") {
		// Check the version of the operating system
		// if !strings.Contains(completeEtcOsReleaseFileContent, "VERSION_ID=\"22.04\"") {
		// 	saveAllErrors("This application is only supported on Ubuntu 22.04.")
		// }
		// Note: Remove in the future build of the app.
		log.Println("Note: This app is only supported on Ubuntu 22.04.")
	} else {
		saveAllErrors("This application is only supported on Ubuntu.")
	}
}

// Append and write to file
func appendAndWriteToFile(path string, content string) {
	filePath, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	_, err = filePath.WriteString(content + "\n")
	if err != nil {
		log.Println(err)
	}
	err = filePath.Close()
	if err != nil {
		log.Println(err)
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
		saveAllErrors("The value for the key '" + key + "' is empty.")
	}
	return false
}

// Move file from one location to another
func moveFile(source string, destination string) {
	// Check if the source file exists
	if fileExists(source) {
		// Check if the destination file exists
		if fileExists(destination) {
			// Remove the destination file
			err := os.Remove(destination)
			if err != nil {
				log.Println(err)
			}
		}
		// Move the file
		err := os.Rename(source, destination)
		if err != nil {
			log.Println(err)
		}
	}
}

// Check if the config changes and exit the application if it does.
func checkConfigChanges() {
	// Get the hash of the config file.
	initialConfigHash := sha256OfFile(applicationConfigFile)
	for {
		// Check if the config file has changed
		if initialConfigHash != sha256OfFile(applicationConfigFile) {
			// Save the error
			saveAllErrors("The config file has changed. Please restart the application.")
		}
		// Sleep for 15 second
		time.Sleep(15 * time.Second)
	}

}

/*
Imports the "os" package which provides the UserHomeDir() function
Defines the currentUserHomeDir() function
Invokes the UserHomeDir() function
Returns the home directory of the current user
Returns -1 if no user home directory is found
*/
func currentUserHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "-1"
	}
	return homeDir
}

// Find the AWS credentials file.
func findAWSCredentialsFile() string {
	// Check if there is a AWS creditentials folder in the home dir.
	if fileExists(currentUserHomeDir() + "/.aws/credentials") {
		return currentUserHomeDir() + "/.aws/credentials"
	}
	return "Unknown"
}

// Parse the AWS credentials file.
func parseAWSCredentialsFile() (string, string) {
	// Define the AWS access key and secret key.
	var awsAccessKey string
	var awsSecretKey string
	// Check if the AWS credentials file exists
	if fileExists(findAWSCredentialsFile()) {
		// Read the AWS credentials file
		awsCredentialsFileContent := readAFileAsString(findAWSCredentialsFile())
		// Split the file into lines.
		awsCredentialsFileContentLines := strings.Split(awsCredentialsFileContent, "\n")
		// Loop through the lines.
		for _, line := range awsCredentialsFileContentLines {
			// Check if the line contains the access key.
			if strings.Contains(line, "aws_access_key_id") {
				// Get the access key.
				awsAccessKey = strings.Split(line, "=")[1]
			}
			// Check if the line contains the secret key.
			if strings.Contains(line, "aws_secret_access_key") {
				// Get the secret key.
				awsSecretKey = strings.Split(line, "=")[1]
			}
		}
		// Remove whitespace from the keys.
		awsAccessKey = strings.TrimSpace(awsAccessKey)
		awsSecretKey = strings.TrimSpace(awsSecretKey)
	}
	if len(awsAccessKey) == 0 {
		saveAllErrors("The AWS access key is empty.")
	}
	if len(awsSecretKey) == 0 {
		saveAllErrors("The AWS secret key is empty.")
	}
	return awsAccessKey, awsSecretKey
}

/*
Checks if the directory exists
If it exists, return true.
If it doesn't, return false.
*/
func directoryExists(path string) bool {
	directory, err := os.Stat(path)
	if err != nil {
		return false
	}
	return directory.IsDir()
}

// Add a key-value pair to the given map.
func addKeyValueToMap(providedMap map[string]bool, key string, value bool) map[string]bool {
	providedMap[key] = value
	return providedMap
}

// Get the value of a key from the given map.
func getValueFromMap(providedMap map[string]bool, key string) bool {
	return providedMap[key]
}

// Get the server info.
func getServerByIndex(config *AutoGenerated, index int) HostStruct {
	switch index {
	case 0:
		return config.Num0
	case 1:
		return config.Num1
	case 2:
		return config.Num2
	default:
		panic("Invalid server index")
	}
}
