package main

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/aler9/gortsplib/v2"
	"github.com/aler9/gortsplib/v2/pkg/format"
	"github.com/aler9/gortsplib/v2/pkg/media"
	"github.com/aler9/gortsplib/v2/pkg/url"
	"github.com/pion/rtp"
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
func saveAllErrors(errors error, path string) {
	filePath, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	log.SetOutput(filePath)
	log.Println(errors)
	err = filePath.Close()
	if err != nil {
		log.Fatalln(err)
	}
}

// exitTheApplication prints the message to the log and exits the application
func exitTheApplication(message string) {
	log.Fatalln(message)
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
	// Connect to the go RTSP server client
	serverConnection := gortsplib.Client{}
	// Start the connection
	err = serverConnection.Start(parsedURL.Scheme, parsedURL.Host)
	if err != nil {
		return false
	}
	// Get the media data from the server connection
	mediaData, baseURL, _, err := serverConnection.Describe(parsedURL)
	if err != nil {
		return false
	}
	// Setup the connection
	err = serverConnection.SetupAll(mediaData, baseURL)
	if err != nil {
		return false
	}
	// Note: This might be wrong and instead use the map to count how many times it happens.
	// I think the last packet loops over and over.
	// List of invalid packets.
	invalidPacketList := []string{
		"&{audio  mediaUUID=59b4572b-6cfa-4424-8bdf-d06b9b31ef8d [MPEG4-audio]}",
	}
	// Counter for the number of packets received
	invalidPacketCounter := 0
	// Invalid return value; kill the connection
	invalidReturnValue := false
	// Get the packet from the server
	serverConnection.OnPacketRTPAny(func(medi *media.Media, forma format.Format, pkt *rtp.Packet) {
		mediaValueAsString := fmt.Sprintf("%v", medi)
		// Add checks to make sure the packets are valid
		for _, invalidPacket := range invalidPacketList {
			if mediaValueAsString == invalidPacket {
				invalidPacketCounter = invalidPacketCounter + 1
				if invalidPacketCounter >= 100 {
					invalidReturnValue = true
					return
				}
			}
		}
	})
	// Play the stream
	_, err = serverConnection.Play(nil)
	if err != nil {
		return false
	}
	// Kill the connection if the return value is invalid
	if invalidReturnValue {
		return false
	}
	// We will watch the connection for 30 seconds
	time.Sleep(30 * time.Second)
	// Close the connection
	defer serverConnection.Close()
	return true
}

// Run this function in the background and check if a given RTSP server is alive
func checkRTSPServerAliveInBackground(rtspURL string) {
	for {
		// Check if the server is alive
		if checkRTSPServerAlive(rtspURL) {
			rtspSeverOneStatus = true
		} else {
			rtspSeverOneStatus = false
		}
		// Check if the server is alive and sleep for 1 minute; else sleep for 30 second
		if rtspSeverOneStatus {
			time.Sleep(1 * time.Minute)
			continue
		} else {
			time.Sleep(30 * time.Second)
			continue
		}
	}
	rtspServerWaitGroup.Done()
}

// Run a command in the system terminal.
func runSystemTerminalCommand(content string) {
	cmd := exec.Command(content)
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
