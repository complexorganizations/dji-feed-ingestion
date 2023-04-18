package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// Check if there is a usb mounted and if there is return it.
func getUSBMountPoint() (string, error) {
	cmd := exec.Command("df", "-h")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "/media") || strings.Contains(line, "/run/media") {
			fields := strings.Fields(line)
			return fields[len(fields)-1], nil
		}
	}
	return "", fmt.Errorf("USB device not found")
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
		log.Fatalln(err)
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
		log.Fatalln(err)
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
func removeFile(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Fatalln(err)
	}
}

// Remove a directory and all its contents.
func removeDirectory(dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		log.Fatalln(err)
	}
}

// Check if a given directory is empty.
func isDirectoryEmpty(path string) bool {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatalln(err)
	}
	return len(files) == 0
}

// Nuke a given directory and all its contents.
func nukeDirectory(path string) {
	if directoryExists(path) {
		if !isDirectoryEmpty(path) {
			removeDirectory(path)
		}
	}
}

// Move a file from one location to another.
func moveFile(source string, destination string) {
	// Get the file name from the source path
	fileName := filepath.Base(source)
	// Move the file to the destination
	log.Println("Moving file: " + source + " to: " + destination + fileName)
	cmd := exec.Command("cp", source, destination+fileName)
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

// Generate a random string of a given length.
func generateRandomString(length int) string {
	randomBytes := make([]byte, length/2)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Fatalln(err)
	}
	return fmt.Sprintf("%x", randomBytes)
}

/* The function takes two parameters: path and permission.
We use os.Mkdir() to create the directory.
If there is an error, we use log.Fatalln() to log the error and then exit the program. */
func createDirectory(path string, permission os.FileMode) {
	err := os.Mkdir(path, permission)
	if err != nil {
		log.Fatalln(err)
	}
}
