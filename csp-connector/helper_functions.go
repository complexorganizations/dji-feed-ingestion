package main

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"

	"github.com/aler9/gortsplib/v2"
	"github.com/aler9/gortsplib/v2/pkg/url"
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

// Check if the given url is valid.
func isUrlValid(uri string) bool {
	_, err := url.ParseRequestURI(uri)
	return err == nil
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
