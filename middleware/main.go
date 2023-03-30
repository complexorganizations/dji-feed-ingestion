package main

import (
	"log"
	"time"
)

func main() {
	time.Sleep(3 * time.Second)
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
				// Upload the feed into Google cloud vertex AI
				forwardDataToGoogleCloudVertexAI(TestJSONValue.Num0.Host, TestJSONValue.Num0.GoogleCloudVertexAiVision.ProjectName, TestJSONValue.Num0.GoogleCloudVertexAiVision.GcpRegion, TestJSONValue.Num0.GoogleCloudVertexAiVision.VertexStreams)
				// Upload the feed into AWS Kinesis Video Streams
				runGstPipeline(TestJSONValue.Num0.Host, TestJSONValue.Num0.AmazonKinesisVideoStreams.KinesisStream, TestJSONValue.Num0.AmazonKinesisVideoStreams.AwsAccessKeyID, TestJSONValue.Num0.AmazonKinesisVideoStreams.AwsSecretAccessKey, TestJSONValue.Num0.AmazonKinesisVideoStreams.AwsDefaultRegion)
				// Docs: Upload mock test data.
				log.Println("We are about to remove a 1 from the counter.")
				// Remove a 1 from the counter when the upload is done
				rtspConnectionCounter = rtspConnectionCounter - 1
			}
		}
		// Sleep for 30 seconds
		time.Sleep(30 * time.Second)
	}
}
