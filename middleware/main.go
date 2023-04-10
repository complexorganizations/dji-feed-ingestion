package main

import (
	"flag"
	"log"
	"os"
	"time"
	"sync"
)

var (
	applicationConfigFile = getCurrentWorkingDirectory() + "config.json"
	applicationLogFile    = getCurrentWorkingDirectory() + "log.txt"
	currentJsonValue      interface{}
	// Note: Future update for now we are using a temp bool var
	// rtspServerStatusChannel = make(map[string]bool)
	rtspServerOneStatus bool
	uploadWaitGroup sync.WaitGroup
	debug               bool
)

// The config file struct for the application to use.
type AutoGenerated struct {
	Num0 struct {
		Host                      string `json:"host"`
		AmazonKinesisVideoStreams struct {
			AccessKeyID     string `json:"access_key_id"`
			SecretAccessKey string `json:"secret_access_key"`
			DefaultRegion   string `json:"default_region"`
			KinesisStream   string `json:"kinesis_stream"`
		} `json:"amazon_kinesis_video_streams"`
		GoogleCloudVertexAiVision struct {
			ProjectName          string `json:"project_name"`
			DefaultRegion        string `json:"default_region"`
			VertexAiVisionStream string `json:"vertex_ai_vision_stream"`
		} `json:"google_cloud_vertex_ai_vision"`
	} `json:"0"`
}

func init() {
	// Validate the operating system
	lockdownToLinuxOperatingSystem()
	// Check if there are any user provided flags in the request.
	if len(os.Args) > 1 {
		// Check if the config path is provided.
		tempConfig := flag.String("config", "config.json", "The location of the config file.")
		tempLog := flag.String("log", "log.txt", "The location of the log file.")
		tempDebug := flag.Bool("debug", false, "Determine if this is a debug run.")
		flag.Parse()
		applicationConfigFile = *tempConfig
		applicationLogFile = *tempLog
		debug = *tempDebug
	} else {
		// if there are no flags provided than we close the application.
		log.Fatalln("Error: No flags provided. Please use -help for more information.")
	}
	// Check if the system has the required tools and is installed in path.
	requiredApplications := []string{"vaictl"}
	// Check if the required application are present in the system
	for _, app := range requiredApplications {
		if !commandExists(app) {
			saveAllErrors("Error: " + app + "is not installed in your system, Please install it and try again.")
		}
	}
	// Check if the config file exists in the current directory
	if !fileExists(applicationConfigFile) {
		// Write a config file in the current directory if it doesn't exist
		writeToFile(applicationConfigFile, []byte(encodeStructToJSON(AutoGenerated{})))
		// Exit the application since the config file was written just now and content will not be in that file.
		saveAllErrors("Error: Just created the default configuration; please edit the configuration and launch the program again.")
	}
	// Check if the file provided has a valid .json extension.
	if !getFileExtension(applicationConfigFile) == ".json" {
		saveAllErrors("Error: The extension of the config file isn't valid.")
	}
	// DEBUG: Print the Hash of the file to change it below.
	// log.Println(sha256OfFile(applicationConfigFile))
	// Hash the file and get the SHA-256 and make sure its not the deafult config.
	if sha256OfFile(applicationConfigFile) == "273dfdef0f9b697b5b76f23e23e17563c9ab56eff100093b5ac1ef411546da15e19c0aae8153c64691a4a86b5db2465bebd6943b863531149b4995a3f55ba0ad" {
		// The file has not been modified
		saveAllErrors("Error: The config file has not been modified, Please modify it and try again.")
	}
	// Check if the config has the correct format and all the info is correct.
	if !jsonValid(readFileAndReturnAsBytes(applicationConfigFile)) {
		saveAllErrors("Error: The config file is not a valid json file.")
	}
	// Read the config file and store it in a variable
	currentJsonValue = unmarshalJSONIntoStruct([]byte(readFileAndReturnAsBytes(applicationConfigFile)), &AutoGenerated{})
	// Make sure the length of the json is not 0
	if len(currentJsonValue.(*AutoGenerated).Num0.Host) == 0 {
		saveAllErrors("Error: The host value is not set in the config file.")
	}
	// Make sure non of the values are deafult; if it is than exit.

	// Validate all the data thats imported in the app; test run the connection if possible.

	// Check if the rtsp server is alive and responding to requests
	go checkRTSPServerAliveInBackground(currentJsonValue.(*AutoGenerated).Num0.Host)
}

func main() {
	// RTSP Connection Counter
	log.Println("We are outside the loop, about to start the loop")
	rtspConnectionCounter := 0
	for {
		log.Println("We are inside the loop")
		// Check the ammount of time the rtsp server has run
		if rtspConnectionCounter == 0 {
			log.Println(rtspServerOneStatus)
			log.Println("We got the status of the rtsp server.")
			// Check if the rtsp server is alive and responding to requests; run the upload in the background
			if rtspServerOneStatus {
				log.Println("The rtsp server is alive and we are about to start the upload.")
				// Add a 1 to the counter
				rtspConnectionCounter = rtspConnectionCounter + 1
				// Add a 1 to the wait group
				uploadWaitGroup.Add(1)
				// Upload the feed into AWS Kinesis Video Streams
				go runGstPipeline(currentJsonValue.(*AutoGenerated).Num0.Host, currentJsonValue.(*AutoGenerated).Num0.AmazonKinesisVideoStreams.KinesisStream, currentJsonValue.(*AutoGenerated).Num0.AmazonKinesisVideoStreams.AccessKeyID, currentJsonValue.(*AutoGenerated).Num0.AmazonKinesisVideoStreams.SecretAccessKey, currentJsonValue.(*AutoGenerated).Num0.AmazonKinesisVideoStreams.DefaultRegion)
				log.Println("We are writing to the aws kinesis video streams.")
				// Note temp work around
				// mv /etc/amazon-kinesis-video-streams-producer-sdk-cpp/build/libgstkvssink.so /etc/amazon-kinesis-video-streams-producer-sdk-cpp/build/libgstkvssink.so.tmp
				// This is a temp work around for the issue with the kinesis video streams, Once your done with the upload you can move the file back to its original location.
				// mv /etc/amazon-kinesis-video-streams-producer-sdk-cpp/build/libgstkvssink.so.tmp /etc/amazon-kinesis-video-streams-producer-sdk-cpp/build/libgstkvssink.so
				// Upload the feed into Google cloud vertex AI
				// go forwardDataToGoogleCloudVertexAI(currentJsonValue.(*AutoGenerated).Num0.Host, currentJsonValue.(*AutoGenerated).Num0.GoogleCloudVertexAiVision.ProjectName, currentJsonValue.(*AutoGenerated).Num0.GoogleCloudVertexAiVision.DefaultRegion, currentJsonValue.(*AutoGenerated).Num0.GoogleCloudVertexAiVision.VertexAiVisionStream)
				// log.Println("We are writing to the google cloud vertex ai.")
				// Wait for the wait group to finish
				uploadWaitGroup.Wait()
				// Docs: Upload mock test data.
				log.Println("We are about to remove a 1 from the counter.")
				// Remove a 1 from the counter when the upload is done
				rtspConnectionCounter = rtspConnectionCounter - 1
			}
		}
		// Sleep for 5 seconds
		time.Sleep(5 * time.Second)
		// debug
		if debug {
			log.Println("Debug session, ending now")
			os.Exit(0)
		}
	}
}
