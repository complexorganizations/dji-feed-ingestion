package main

import (
	"context"
	"flag"
	"log"
	"os"
	"reflect"
	"strconv"
	"sync"
	"time"
)

var (
	// Application config file
	applicationConfigFile = getCurrentWorkingDirectory() + "config.json"
	// Application log file
	applicationLogFile         = getCurrentWorkingDirectory() + "log.txt"
	currentJsonValue           interface{}
	rtspServerStatusChannel    = make(map[string]bool)
	rtspServerStreamingChannel = make(map[string]bool)
	rtspServerPacketChannel    = make(map[string]bool)
	debug                      bool
	awsKVS                     bool
	awsIVS                     bool
	gcp                        bool
	yt                         bool
	twitch                     bool
	facebook                   bool
	rtmp                       bool
	// Values for the aws file path stuff;
	amazonKinesisVideoStreamPath      = "/etc/amazon-kinesis-video-streams-producer-sdk-cpp/"
	amazonKinesisVideoStreamBuildPath = amazonKinesisVideoStreamPath + "build/"
	// Note: There is an issue with vertex ai vision; where a file needs to be moved. https://github.com/google/visionai/issues/6
	amazonKinesisDefaultPath = amazonKinesisVideoStreamBuildPath + "libgstkvssink.so"
	amazonKinesisTempPath    = amazonKinesisDefaultPath + ".tmp"
	// Google Cloud Credentials
	googleCloudCredentials = currentUserHomeDir() + "/.config/gcloud/application_default_credentials.json"
	// Number of clients allowed on the system.
	numberOfClientsAllowed int
)

// The config file struct for the application to use.
type AutoGenerated struct {
	Num0  HostStruct `json:"0"`
	Num1  HostStruct `json:"1"`
	Num2  HostStruct `json:"2"`
	Num3  HostStruct `json:"3"`
	Num4  HostStruct `json:"4"`
	Num5  HostStruct `json:"5"`
	Num6  HostStruct `json:"6"`
	Num7  HostStruct `json:"7"`
	Num8  HostStruct `json:"8"`
	Num9  HostStruct `json:"9"`
	Num10 HostStruct `json:"10"`
	Num11 HostStruct `json:"11"`
	Num12 HostStruct `json:"12"`
	Num13 HostStruct `json:"13"`
	Num14 HostStruct `json:"14"`
	Num15 HostStruct `json:"15"`
	Num16 HostStruct `json:"16"`
	Num17 HostStruct `json:"17"`
	Num18 HostStruct `json:"18"`
	Num19 HostStruct `json:"19"`
	Num20 HostStruct `json:"20"`
	Num21 HostStruct `json:"21"`
	Num22 HostStruct `json:"22"`
	Num23 HostStruct `json:"23"`
	Num24 HostStruct `json:"24"`
	Num25 HostStruct `json:"25"`
	Num26 HostStruct `json:"26"`
	Num27 HostStruct `json:"27"`
	Num28 HostStruct `json:"28"`
	Num29 HostStruct `json:"29"`
	Num30 HostStruct `json:"30"`
	Num31 HostStruct `json:"31"`
	Num32 HostStruct `json:"32"`
}

type HostStruct struct {
	Host                          string                        `json:"host"`
	AmazonKinesisVideoStreams     AmazonKinesisVideoStreams     `json:"amazon_kinesis_video_streams"`
	AmazonInteractiveVideoService AmazonInteractiveVideoService `json:"amazon_interactive_video_service"`
	GoogleCloudVertexAiVision     GoogleCloudVertexAiVision     `json:"google_cloud_vertex_ai_vision"`
	YoutubeLiveStream             YoutubeLiveStream             `json:"youtube_live_stream"`
	TwitchLiveStream              TwitchLiveStream              `json:"twitch_live_stream"`
	FacebookLiveStream            FacebookLiveStream            `json:"facebook_live_stream"`
	RtmpServer                    RtmpServer                    `json:"rtmp_server"`
}

type AmazonKinesisVideoStreams struct {
	DefaultRegion string `json:"default_region"`
	KinesisStream string `json:"kinesis_stream"`
}

type AmazonInteractiveVideoService struct {
	IvsStream string `json:"ivs_stream"`
}

type GoogleCloudVertexAiVision struct {
	ProjectName          string `json:"project_name"`
	DefaultRegion        string `json:"default_region"`
	VertexAiVisionStream string `json:"vertex_ai_vision_stream"`
}

type YoutubeLiveStream struct {
	StreamKey string `json:"stream_key"`
}

type TwitchLiveStream struct {
	StreamKey string `json:"stream_key"`
}

type FacebookLiveStream struct {
	StreamKey string `json:"stream_key"`
}

type RtmpServer struct {
	ConnectionString string `json:"connection_string"`
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
		tempAWSKVS := flag.Bool("aws_kvs", false, "Determine if this is a AWS run.")
		tempAWSIVS := flag.Bool("aws_ivs", false, "Determine if this is a AWS run.")
		tempGCP := flag.Bool("gcp", false, "Determine if this is a GCP run.")
		tempYT := flag.Bool("yt", false, "Determine if this is a YT run.")
		tempTwitch := flag.Bool("twitch", false, "Determine if this is a Twitch run.")
		tempFacebook := flag.Bool("facebook", false, "Determine if this is a Facebook run.")
		tempRTMP := flag.Bool("rtmp", false, "Determine if this is a Any RTMP run.")
		flag.Parse()
		// Set the values to the global variables.
		applicationConfigFile = *tempConfig
		applicationLogFile = *tempLog
		debug = *tempDebug
		awsKVS = *tempAWSKVS
		awsIVS = *tempAWSIVS
		gcp = *tempGCP
		yt = *tempYT
		twitch = *tempTwitch
		facebook = *tempFacebook
		rtmp = *tempRTMP
	} else {
		// if there are no flags provided than we close the application.
		log.Fatalln("Error: No flags provided. Please use -help for more information.")
	}
	// Only run one of the three options.
	if awsKVS && awsIVS && gcp && yt && twitch && facebook && rtmp {
		log.Fatalln("Error: You can only run one of the -help options.")
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
	if sha256OfFile(applicationConfigFile) == "7f6b889e1e2f5dbbbf4c44b36c875de4a6ed6310f3533ebdaefbeed5e5fb10107f8526b1b068feadfab5966f04d398314f444b749c2e5e7c42171c22800b5753" {
		// The file has not been modified
		saveAllErrors("Error: The config file has not been modified, Please modify it and try again.")
	}
	// Check if the config has the correct format and all the info is correct.
	if !jsonValid(readFileAndReturnAsBytes(applicationConfigFile)) {
		saveAllErrors("Error: The config file is not a valid json file.")
	}
	// Read the config file and store it in a variable
	currentJsonValue = unmarshalJSONIntoStruct([]byte(readFileAndReturnAsBytes(applicationConfigFile)), &AutoGenerated{})
	// Get the lenth of the struct; this will be the number of servers.
	numberOfHosts := reflect.TypeOf(AutoGenerated{})
	numberOfClientsAllowed = numberOfHosts.NumField()
	// Checks how many hosts are in the config file and than determines if the app will allow it.
	if countHosts() >= numberOfClientsAllowed {
		log.Fatalln("Warning: The number of servers in the config file is more than the number of servers in the struct.")
	}
	// Get the real number of servers in the config file.
	numberOfClientsAllowed = countHosts()
	// RTSP Server Counter Map.
	for i := 0; i < numberOfClientsAllowed; i++ {
		server := getServerByIndex(currentJsonValue.(*AutoGenerated), i)
		// Validate the length of the config file
		validateJSONLength("Host", server.Host)
		// AWS KVS
		validateJSONLength("Amazon Kinesis Region", server.AmazonKinesisVideoStreams.DefaultRegion)
		validateJSONLength("Amazon Kinesis Stream Name", server.AmazonKinesisVideoStreams.KinesisStream)
		// AWS IVS
		validateJSONLength("Amazon IVS Stream Name", server.AmazonInteractiveVideoService.IvsStream)
		// GCP
		validateJSONLength("Google Project Name", server.GoogleCloudVertexAiVision.ProjectName)
		validateJSONLength("Google Default Region", server.GoogleCloudVertexAiVision.DefaultRegion)
		validateJSONLength("Google Vertex AI Vision Stream", server.GoogleCloudVertexAiVision.VertexAiVisionStream)
		// YouTube
		validateJSONLength("YouTube Stream Name", server.YoutubeLiveStream.StreamKey)
		// Twitch
		validateJSONLength("Twitch Stream Name", server.TwitchLiveStream.StreamKey)
		// Facebook
		validateJSONLength("Facebook Stream Name", server.FacebookLiveStream.StreamKey)
		// RTMP
		validateJSONLength("RTMP Stream Name", server.RtmpServer.ConnectionString)
		// Check if the rtsp server is alive and responding to requests
		go checkRTSPServerAliveInBackground(server.Host)
	}
	// Note: This is a temp location for this and other location will be better for this.
	go checkConfigChanges()
	// Validate the content of the config file (API Keys, etc.)
}

func main() {
	// Setup the variables for aws.
	var accessKey string
	var secretKey string
	if awsIVS || awsKVS {
		// Get the AWS Credentials
		accessKey, secretKey = parseAWSCredentialsFile()
	} else if gcp {
		// Get the Google Cloud Credentials
		validateGoogleCloudCLI()
	}
	// Create a wait group for the upload
	var uploadWaitGroup sync.WaitGroup
	// Create a counter map for the RTSP Server
	var rtspServerRunCounter = make(map[string]int)
	// Var counter
	var counter int
	// Create a map for the cancel functions
	var cancelFuncs = make(map[string]context.CancelFunc)
	for {
		for i := 0; i < numberOfClientsAllowed; i++ {
			// Start the cancel and the context for the RTSP Server
			ctx, cancel := context.WithCancel(context.Background())
			// Get the server from the config file
			server := getServerByIndex(currentJsonValue.(*AutoGenerated), i)
			log.Println("#1:" + server.Host + strconv.FormatBool(getValueFromMap(rtspServerStatusChannel, server.Host)))
			// Make sure the server isnt already streaming the feed.
			if getValueFromMap(rtspServerStreamingChannel, server.Host) == false {
				// Make sure you only run the server once.
				if rtspServerRunCounter[server.Host] == 0 {
					// Prevent the server from running again if it is already running
					rtspServerRunCounter[server.Host] = 1
					log.Println("#2:" + server.Host + strconv.FormatBool(getValueFromMap(rtspServerStatusChannel, server.Host)))
					// Check if the server is alive and responding to requests
					if getValueFromMap(rtspServerStatusChannel, server.Host) {
						// Add key-value pair to the map
						cancelFuncs[server.Host] = cancel
						// Check the packet length in the background
						go checkRTSPStreamPacketConnectionInLoop(server.Host)
						// Increment the counter
						counter = counter + 1
						// Log the counter
						log.Println("#3:" + server.Host + strconv.FormatBool(getValueFromMap(rtspServerStatusChannel, server.Host)))
						// Add one to the wait group
						uploadWaitGroup.Add(1)
						if awsKVS {
							go forwardDataToAmazonKinesisStreams(server.Host, server.AmazonKinesisVideoStreams.KinesisStream, accessKey, secretKey, server.AmazonKinesisVideoStreams.DefaultRegion, &uploadWaitGroup, ctx)
						} else if gcp {
							go forwardDataToGoogleCloudVertexAI(server.Host, server.GoogleCloudVertexAiVision.ProjectName, server.GoogleCloudVertexAiVision.DefaultRegion, server.GoogleCloudVertexAiVision.VertexAiVisionStream, &uploadWaitGroup, ctx)
						} else if awsIVS {
							go forwardDataToAmazonIVS(server.Host, server.AmazonInteractiveVideoService.IvsStream, accessKey, secretKey, &uploadWaitGroup, ctx)
						} else if yt {
							go forwardDataToYoutubeLive(server.Host, server.YoutubeLiveStream.StreamKey, &uploadWaitGroup, ctx)
						} else if twitch {
							go forwardDataToTwitch(server.Host, server.TwitchLiveStream.StreamKey, &uploadWaitGroup, ctx)
						} else if facebook {
							go forwardDataToFacebookLive(server.Host, server.FacebookLiveStream.StreamKey, &uploadWaitGroup, ctx)
						} else if rtmp {
							go forwardDataToAnyRTMP(server.Host, server.RtmpServer.ConnectionString, &uploadWaitGroup, ctx)
						}
					}
					rtspServerRunCounter[server.Host] = 0
				}
			}
			// Task: Cancel the context if the server is not alive.
			// Check if there is a upload in progress.
			if getValueFromMap(rtspServerStreamingChannel, server.Host) == true {
				// Check if the server is alive and responding to requests
				if getValueFromMap(rtspServerStatusChannel, server.Host) == false {
					// Cancel the context
					cancelFuncs[server.Host]()
				}
				// Check if the packet connection is alive
				if getValueFromMap(rtspServerPacketChannel, server.Host) == false {
					// Cancel the context
					cancelFuncs[server.Host]()
				}
			}
			// Debug
			log.Println("Map", rtspServerStreamingChannel)
		}
		// This sleep determins how often the program checks if the RTSP server is alive and streaming.
		time.Sleep(5 * time.Second)
		// The counter for how many streams are being uploaded.
		log.Println("Counter: " + strconv.Itoa(counter))
		// End if debug
		if debug {
			break
		}
	}
	// Wait for the upload to finish
	uploadWaitGroup.Wait()
}
