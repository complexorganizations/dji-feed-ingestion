package main

import (
	"flag"
	"log"
	"os"
	"sync"
	"time"
)

var (
	// Application config file
	applicationConfigFile = getCurrentWorkingDirectory() + "config.json"
	// Application log file
	applicationLogFile      = getCurrentWorkingDirectory() + "log.txt"
	currentJsonValue        interface{}
	rtspServerStatusChannel = make(map[string]bool)
	debug                   bool
	aws                     bool
	gcp                     bool
	// Values for the aws file path stuff;
	amazonKinesisVideoStreamPath      = "/etc/amazon-kinesis-video-streams-producer-sdk-cpp"
	amazonKinesisVideoStreamBuildPath = amazonKinesisVideoStreamPath + "/build"
	// This is the issue with in google to fix this stuff. /// https://github.com/google/visionai/issues/6
	amazonKinesisDefaultPath = amazonKinesisVideoStreamBuildPath + "/libgstkvssink.so"
	amazonKinesisTempPath    = amazonKinesisDefaultPath + ".tmp"
)

// The config file struct for the application to use.
type AutoGenerated struct {
	Num0 HostStruct `json:"0"`
	Num1 HostStruct `json:"1"`
	Num2 HostStruct `json:"2"`
	Num3 HostStruct `json:"3"`
	Num4 HostStruct `json:"4"`
	Num5 HostStruct `json:"5"`
}

type HostStruct struct {
	Host                      string                    `json:"host"`
	AmazonKinesisVideoStreams AmazonKinesisVideoStreams `json:"amazon_kinesis_video_streams"`
	GoogleCloudVertexAiVision GoogleCloudVertexAiVision `json:"google_cloud_vertex_ai_vision"`
}

type AmazonKinesisVideoStreams struct {
	DefaultRegion string `json:"default_region"`
	KinesisStream string `json:"kinesis_stream"`
}

type GoogleCloudVertexAiVision struct {
	ProjectName          string `json:"project_name"`
	DefaultRegion        string `json:"default_region"`
	VertexAiVisionStream string `json:"vertex_ai_vision_stream"`
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
		tempAWS := flag.Bool("aws", false, "Determine if this is a AWS run.")
		tempGCP := flag.Bool("gcp", false, "Determine if this is a GCP run.")
		flag.Parse()
		applicationConfigFile = *tempConfig
		applicationLogFile = *tempLog
		debug = *tempDebug
		aws = *tempAWS
		gcp = *tempGCP
	} else {
		// if there are no flags provided than we close the application.
		log.Fatalln("Error: No flags provided. Please use -help for more information.")
	}
	// Both AWS and GCP can't be true at the same time.
	if aws && gcp {
		saveAllErrors("Error: Both AWS and GCP can't be true at the same time.")
	}
	if !aws && !gcp {
		saveAllErrors("Error: Both AWS and GCP can't be false at the same time.")
	}
	// Check if the system has the required tools and is installed in path.
	requiredApplications := []string{
		"vaictl",
		"gst-launch-1.0",
		"ffmpeg",
		"aws",
		"gcloud",
	}
	// Check if the required application are present in the system
	for _, app := range requiredApplications {
		if !commandExists(app) {
			saveAllErrors("Error: " + app + "is not installed in your system, Please install it and try again.")
		}
	}
	// Check the directory structure for the application
	// Check if the amazon-kinesis-video-streams-producer-sdk-cpp is present in the system
	if !directoryExists(amazonKinesisVideoStreamPath) {
		saveAllErrors("Error: The amazon-kinesis-video-streams-producer-sdk-cpp is not present in the system.")
	}
	// Check if the amazon-kinesis-video-streams-producer-sdk-cpp/build is present in the system
	if !directoryExists(amazonKinesisVideoStreamBuildPath) {
		saveAllErrors("Error: The amazon-kinesis-video-streams-producer-sdk-cpp/build is not present in the system.")
	}
	// Check if the amazon-kinesis-video-streams-producer-sdk-cpp/build/libgstkvssink.so is present in the system
	if !fileExists(amazonKinesisDefaultPath) && !fileExists(amazonKinesisTempPath) {
		saveAllErrors("Error: The amazon-kinesis-video-streams-producer-sdk-cpp/build/libgstkvssink.so is not present in the system.")
	}
	// Check if the config file exists in the current directory
	if !fileExists(applicationConfigFile) {
		// Write a config file in the current directory if it doesn't exist
		writeToFile(applicationConfigFile, []byte(encodeStructToJSON(AutoGenerated{})))
		// Exit the application since the config file was written just now and content will not be in that file.
		saveAllErrors("Error: Just created the default configuration; please edit the configuration and launch the program again.")
	}
	// Check if the file provided has a valid .json extension.
	if getFileExtension(applicationConfigFile) != ".json" {
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
	log.Println(currentJsonValue)

	// Validate the first level of the config file
	// Validate the length of the config file
	validateJSONLength("Host", currentJsonValue.(*AutoGenerated).Num0.Host)
	validateJSONLength("DefaultRegion", currentJsonValue.(*AutoGenerated).Num0.AmazonKinesisVideoStreams.DefaultRegion)
	validateJSONLength("KinesisStream", currentJsonValue.(*AutoGenerated).Num0.AmazonKinesisVideoStreams.KinesisStream)
	validateJSONLength("ProjectName", currentJsonValue.(*AutoGenerated).Num0.GoogleCloudVertexAiVision.ProjectName)
	validateJSONLength("DefaultRegion", currentJsonValue.(*AutoGenerated).Num0.GoogleCloudVertexAiVision.DefaultRegion)
	validateJSONLength("VertexAiVisionStream", currentJsonValue.(*AutoGenerated).Num0.GoogleCloudVertexAiVision.VertexAiVisionStream)
	// Validate the second level of the config file
	validateJSONLength("AccessKey", currentJsonValue.(*AutoGenerated).Num1.Host)
	validateJSONLength("SecretKey", currentJsonValue.(*AutoGenerated).Num1.AmazonKinesisVideoStreams.DefaultRegion)
	validateJSONLength("KinesisStream", currentJsonValue.(*AutoGenerated).Num1.AmazonKinesisVideoStreams.KinesisStream)
	validateJSONLength("ProjectName", currentJsonValue.(*AutoGenerated).Num1.GoogleCloudVertexAiVision.ProjectName)
	validateJSONLength("DefaultRegion", currentJsonValue.(*AutoGenerated).Num1.GoogleCloudVertexAiVision.DefaultRegion)
	validateJSONLength("VertexAiVisionStream", currentJsonValue.(*AutoGenerated).Num1.GoogleCloudVertexAiVision.VertexAiVisionStream)
	// Validate the third level of the config file
	validateJSONLength("AccessKey", currentJsonValue.(*AutoGenerated).Num2.Host)
	validateJSONLength("SecretKey", currentJsonValue.(*AutoGenerated).Num2.AmazonKinesisVideoStreams.DefaultRegion)
	validateJSONLength("KinesisStream", currentJsonValue.(*AutoGenerated).Num2.AmazonKinesisVideoStreams.KinesisStream)
	validateJSONLength("ProjectName", currentJsonValue.(*AutoGenerated).Num2.GoogleCloudVertexAiVision.ProjectName)
	validateJSONLength("DefaultRegion", currentJsonValue.(*AutoGenerated).Num2.GoogleCloudVertexAiVision.DefaultRegion)
	validateJSONLength("VertexAiVisionStream", currentJsonValue.(*AutoGenerated).Num2.GoogleCloudVertexAiVision.VertexAiVisionStream)
	// Validate the fourth level of the config file
	validateJSONLength("AccessKey", currentJsonValue.(*AutoGenerated).Num3.Host)
	validateJSONLength("SecretKey", currentJsonValue.(*AutoGenerated).Num3.AmazonKinesisVideoStreams.DefaultRegion)
	validateJSONLength("KinesisStream", currentJsonValue.(*AutoGenerated).Num3.AmazonKinesisVideoStreams.KinesisStream)
	validateJSONLength("ProjectName", currentJsonValue.(*AutoGenerated).Num3.GoogleCloudVertexAiVision.ProjectName)
	validateJSONLength("DefaultRegion", currentJsonValue.(*AutoGenerated).Num3.GoogleCloudVertexAiVision.DefaultRegion)
	validateJSONLength("VertexAiVisionStream", currentJsonValue.(*AutoGenerated).Num3.GoogleCloudVertexAiVision.VertexAiVisionStream)
	// Validate the fifth level of the config file
	validateJSONLength("AccessKey", currentJsonValue.(*AutoGenerated).Num4.Host)
	validateJSONLength("SecretKey", currentJsonValue.(*AutoGenerated).Num4.AmazonKinesisVideoStreams.DefaultRegion)
	validateJSONLength("KinesisStream", currentJsonValue.(*AutoGenerated).Num4.AmazonKinesisVideoStreams.KinesisStream)
	validateJSONLength("ProjectName", currentJsonValue.(*AutoGenerated).Num4.GoogleCloudVertexAiVision.ProjectName)
	validateJSONLength("DefaultRegion", currentJsonValue.(*AutoGenerated).Num4.GoogleCloudVertexAiVision.DefaultRegion)
	validateJSONLength("VertexAiVisionStream", currentJsonValue.(*AutoGenerated).Num4.GoogleCloudVertexAiVision.VertexAiVisionStream)
	// Validate the sixth level of the config file
	validateJSONLength("AccessKey", currentJsonValue.(*AutoGenerated).Num5.Host)
	validateJSONLength("SecretKey", currentJsonValue.(*AutoGenerated).Num5.AmazonKinesisVideoStreams.DefaultRegion)
	validateJSONLength("KinesisStream", currentJsonValue.(*AutoGenerated).Num5.AmazonKinesisVideoStreams.KinesisStream)
	validateJSONLength("ProjectName", currentJsonValue.(*AutoGenerated).Num5.GoogleCloudVertexAiVision.ProjectName)
	validateJSONLength("DefaultRegion", currentJsonValue.(*AutoGenerated).Num5.GoogleCloudVertexAiVision.DefaultRegion)
	validateJSONLength("VertexAiVisionStream", currentJsonValue.(*AutoGenerated).Num5.GoogleCloudVertexAiVision.VertexAiVisionStream)
	// Note: This is a temp location for this and other location will be better for this.
	go checkConfigChanges()
	// Validate the content of the config file (API Keys, etc.)

	// Check if the rtsp server is alive and responding to requests
	// Number of servers
	const numServers = 6
	// RTSP Server Counter Map.
	for i := 0; i < numServers; i++ {
		server := getServerByIndex(currentJsonValue.(*AutoGenerated), i)
		host := server.Host
		log.Println("Checking server: " + host)
		// Check if the server is alive
		go checkRTSPServerAliveInBackground(host)
	}
}

func main() {
	// Get the AWS Credentials
	accessKey, secretKey := parseAWSCredentialsFile()
	// Create a wait group for the upload
	var uploadWaitGroup sync.WaitGroup
	// Create a counter map for the RTSP Server
	var rtspServerRunCounter = make(map[string]int)
	// Number of servers
	const numServers = 6

	// RTSP Server Counter Map.
	for {
		for i := 0; i < numServers; i++ {
			server := getServerByIndex(currentJsonValue.(*AutoGenerated), i)
			host := server.Host
			log.Println("Checking server in second loop: " + host)

			if rtspServerRunCounter[host] == 0 {
				// Add 1 to the counter
				rtspServerRunCounter[host] = 1

				if getValueFromMap(rtspServerStatusChannel, host) {
					uploadWaitGroup.Add(1)

					if aws {
						go runGstPipeline(host, server.AmazonKinesisVideoStreams.KinesisStream, accessKey, secretKey, server.AmazonKinesisVideoStreams.DefaultRegion, &uploadWaitGroup)
					} else if gcp {
						go forwardDataToGoogleCloudVertexAI(host, server.GoogleCloudVertexAiVision.ProjectName, server.GoogleCloudVertexAiVision.DefaultRegion, server.GoogleCloudVertexAiVision.VertexAiVisionStream, &uploadWaitGroup)
					}

					rtspServerRunCounter[host] = 0
				}
			}
		}
		// Wait for the wait group to finish
		uploadWaitGroup.Wait()
		// Sleep for 30 seconds
		time.Sleep(30 * time.Second)
		// End if debug
		if debug {
			break
		}
	}
}

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
	case 5:
		return config.Num5
	default:
		panic("Invalid server index")
	}
}
