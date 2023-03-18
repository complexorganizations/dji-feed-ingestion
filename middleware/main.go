package main

import (
	"log"
	// "time"
)

func main() {
	// RTSP Connection Counter
	rtspConnectionCounter := 0
	for {
		log.Println("We are here")
		// Check the ammount of time the rtsp server has run
		if rtspConnectionCounter == 0 {
			log.Println(<-rtspServerStatusChannel)
			// Check if the rtsp server is alive and responding to requests; run the upload in the background
			if <-rtspServerStatusChannel {
				// Add a 1 to the counter
				rtspConnectionCounter = rtspConnectionCounter + 1
				log.Println("We are here 2")
				// Forward the data to google cloud vertex ai
				forwardDataToGoogleCloudVertexAI(TestJSONValue.Num0.Host, TestJSONValue.Num0.GoogleCloudVertexAiVision.ProjectName, TestJSONValue.Num0.GoogleCloudVertexAiVision.GcpRegion, TestJSONValue.Num0.GoogleCloudVertexAiVision.VertexStreams)
				// Remove a 1 from the counter; the rtsp server has stopped
				rtspConnectionCounter = rtspConnectionCounter - 1
				// vaictl -p github-code-snippets -l us-central1 -c application-cluster-0 --service-endpoint visionai.googleapis.com send rtsp to streams dji-stream-0 --rtsp-uri rtsp://Administrator:Password@localhost:8554/test_0
			}
		}
		// Sleep for 30 seconds
		// time.Sleep(30 * time.Second)
	}
}
