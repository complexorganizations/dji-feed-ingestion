package main

import (
	"context"
	"crypto/rand"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/bluenviron/gortsplib/v3"
	"github.com/bluenviron/gortsplib/v3/pkg/formats"
	"github.com/bluenviron/gortsplib/v3/pkg/media"
	"github.com/bluenviron/gortsplib/v3/pkg/url"
	"github.com/pion/rtp"
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
			go addKeyValueToMap(rtspServerStatusChannel, rtspURL, true)
		} else {
			go addKeyValueToMap(rtspServerStatusChannel, rtspURL, false)
		}
		// If the server is alive check every 1 Minute or check every 3 seconds.
		if checkRTSPServerAlive(rtspURL) {
			time.Sleep(1 * time.Minute)
		} else {
			time.Sleep(3 * time.Second)
		}
	}
}

// Forward data to google cloud vertex AI.
func forwardDataToGoogleCloudVertexAI(host string, projectName string, gcpRegion string, vertexStreams string, forwardingWaitGroup *sync.WaitGroup, ctx context.Context) {
	// Set the rtspServerStreamingChannel to true
	go addKeyValueToMap(rtspServerStreamingChannel, host, true)
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", googleCloudCredentials)
	// Run the command to forward the data to vertex AI
	cmd := exec.CommandContext(ctx, "vaictl", "-p", projectName, "-l", gcpRegion, "-c", "application-cluster-0", "--service-endpoint", "visionai.googleapis.com", "send", "rtsp", "to", "streams", vertexStreams, "--rtsp-uri", host)
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
	// Set the rtspServerStreamingChannel to false
	go addKeyValueToMap(rtspServerStreamingChannel, host, false)
	// Done forwarding
	defer forwardingWaitGroup.Done()
}

// Forward data to AWS Kinesis Video Streams using gstreamer.
func forwardDataToAmazonKinesisStreams(host string, streamName string, accessKey string, secretKey string, awsRegion string, forwardingWaitGroup *sync.WaitGroup, ctx context.Context) {
	// Set the rtspServerStreamingChannel to true
	go addKeyValueToMap(rtspServerStreamingChannel, host, true)
	/*
		// NOTE: THIS IS METHORD 0
		os.Setenv("AWS_ACCESS_KEY_ID", accessKey)
		os.Setenv("AWS_SECRET_ACCESS_KEY", secretKey)
		os.Setenv("AWS_DEFAULT_REGION", awsRegion)
		cmd := exec.Command("./kvs_gstreamer_sample", streamName, host)
		cmd.Dir = amazonKinesisVideoStreamBuildPath
		err := cmd.Run()
		if err != nil {
			log.Println(err)
		}
	*/
	// NOTE: THIS IS METHORD 1
	// Run the gstreamer command to forward the data to AWS Kinesis Video Streams
	os.Setenv("GST_PLUGIN_PATH", "/etc/amazon-kinesis-video-streams-producer-sdk-cpp/build:$GST_PLUGIN_PATH")
	os.Setenv("LD_LIBRARY_PATH", "/etc/amazon-kinesis-video-streams-producer-sdk-cpp/open-source/local/lib:$LD_LIBRARY_PATH")
	cmd := exec.CommandContext(ctx, "gst-launch-1.0", "rtspsrc", "location="+host, "!", "rtph264depay", "!", "h264parse", "!", "kvssink", "stream-name="+streamName, "access-key="+accessKey, "secret-key="+secretKey, "aws-region="+awsRegion)
	/* DEBUG:
	// Redirect the command's stdout and stderr to the current process's stdout
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	*/
	// Run the command
	err := cmd.Run()
	// Check if there was an error
	if err != nil {
		log.Println(err)
	}
	// Set the rtspServerStreamingChannel to false
	go addKeyValueToMap(rtspServerStreamingChannel, host, false)
	// Close the channel.
	defer forwardingWaitGroup.Done()
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
		if !strings.Contains(completeEtcOsReleaseFileContent, "VERSION_ID=\"22") {
			saveAllErrors("This application is only supported on Ubuntu 22")
		}
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
	// Check if the length of the value is 0 or if its just whitespaces.
	if len(value) == 0 || len(strings.TrimSpace(value)) == 0 {
		log.Println("The value for the key '" + key + "' is empty.")
		return false
	}
	// Check if both key and value have a length.
	if len(key) >= 1 && len(value) >= 1 {
		return true
	}
	// If the values are null its okay.
	if value == "null" {
		return true
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
			saveAllErrors("The config was changed and we need to re-run the app to do the checks again.")
		}
		// Update the hash of the config file.
		initialConfigHash = sha256OfFile(applicationConfigFile)
		// Sleep for 1 Minute
		time.Sleep(1 * time.Minute)
	}
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
	// Check if aws sts get-caller-identity is installed
	if validateAWSSTSCallerIdentityCommand() {
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
		} else {
			saveAllErrors("Error: Missing the AWS TS Caller ID.")
		}
	}
	// Check if the AWS access key is empty.
	if len(awsAccessKey) == 0 {
		saveAllErrors("Error: The AWS Access Key is missing.")
	}
	// Check if the AWS secret key is empty.
	if len(awsSecretKey) == 0 {
		saveAllErrors("Error: The AWS secret Key is missing.")
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
	mutex.Lock()
	providedMap[key] = value
	mutex.Unlock()
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
	case 3:
		return config.Num3
	case 4:
		return config.Num4
	default:
		panic("Invalid server index")
	}
}

// Count how many hosts are in the config file.
func countHosts() int {
	// Read the config file.
	configContent := readAFileAsString(applicationConfigFile)
	// Check how many times the word "host" appears in the config file.
	return strings.Count(configContent, "amazon_kinesis_video_streams")
}

// Validate the AWS STS GetCallerIdentity command.
func validateAWSSTSCallerIdentityCommand() bool {
	cmd := exec.Command("aws", "sts", "get-caller-identity")
	out, err := cmd.Output()
	if err != nil {
		log.Println(err)
	}
	// Check if the output contains the word "arn"
	return strings.Contains(string(out), "arn")
}

// Validate that Google cloud cli is authenticated.
func validateGoogleCloudCLI() {
	cmd := exec.Command("gcloud", "auth", "list", "--filter=status:ACTIVE", "--format=value(account)")
	out, err := cmd.Output()
	if err != nil {
		log.Println(err)
	}
	// Exit the app if google cloud creds are there.
	if len(strings.TrimSpace(string(out))) < 5 {
		saveAllErrors("Error: Didn't find any account via the google cloud cli.")
	}
	// Google Cloud Credentials File.
	if !fileExists(googleCloudCredentials) {
		saveAllErrors("Error: Didn't find any google cloud file at " + googleCloudCredentials)
	}
}

// Get a random element from the slice and return the element.
func randomElementFromSlice(slice []string) string {
	someRandomNumber, err := rand.Int(rand.Reader, big.NewInt(int64(len(slice))))
	if err != nil {
		log.Println(err)
	}
	return slice[someRandomNumber.Int64()]
}

// Check the length of a RTSP stream packet and determine if the client is still connected.
func checkRTSPStreamPacketConnection(host string) bool {
	// Create a new client
	rtspClient := gortsplib.Client{}
	// Parse the URL
	url, err := url.Parse(host)
	// Log any errors
	if err != nil {
		log.Println(err)
	}
	// Connect to the server
	err = rtspClient.Start(url.Scheme, url.Host)
	// Log any errors
	if err != nil {
		log.Println(err)
	}
	// Find published medias
	medias, baseURL, _, err := rtspClient.Describe(url)
	// Log any errors
	if err != nil {
		log.Println(err)
	}
	// Setup all medias
	err = rtspClient.SetupAll(medias, baseURL)
	// Log any errors
	if err != nil {
		log.Println(err)
	}
	// Create two variables to store the payload length
	var lessThanFiveHundred int
	var greaterThanFiveHundred int
	// Called when a RTP packet arrives
	rtspClient.OnPacketRTPAny(func(medi *media.Media, forma formats.Format, pkt *rtp.Packet) {
		if len(pkt.Payload) < 500 {
			lessThanFiveHundred = lessThanFiveHundred + 1
		} else if len(pkt.Payload) > 500 {
			greaterThanFiveHundred = greaterThanFiveHundred + 1
		}
	})
	// Start playing
	_, err = rtspClient.Play(nil)
	// Log any errors
	if err != nil {
		log.Println(err)
	}
	// Wait for 5 seconds and then close the connection
	time.Sleep(5 * time.Second)
	// Close the connection
	rtspClient.Close()
	// Return true if there are more packets with a payload length than 500
	return greaterThanFiveHundred > lessThanFiveHundred
}

// Check the packet length of a RTSP stream in a loop in the background.
func checkRTSPStreamPacketConnectionInLoop(host string) {
	for {
		// Check if the server is online
		if getValueFromMap(rtspServerStatusChannel, host) {
			// Check if the packet are still being sent
			if checkRTSPStreamPacketConnection(host) {
				go addKeyValueToMap(rtspServerPacketChannel, host, true)
			} else {
				go addKeyValueToMap(rtspServerPacketChannel, host, false)
			}
		} else {
			go addKeyValueToMap(rtspServerPacketChannel, host, false)
		}
		// If the server is alive, check every 1 Minute or check every 3 seconds.
		if getValueFromMap(rtspServerPacketChannel, host) {
			time.Sleep(1 * time.Minute)
		} else {
			time.Sleep(3 * time.Second)
		}
	}
}
