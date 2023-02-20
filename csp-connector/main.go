package main

func main() {
	// RTSP Connection Counter
	rtspConnectionCounter := 0
	for {
		// Check the ammount of time the rtsp server has run
		if rtspConnectionCounter == 0 {
			// Check if the rtsp server is alive and responding to requests; run the upload in the background
			if rtspSeverOneStatus {
				// Add a 1 to the counter
				rtspConnectionCounter = rtspConnectionCounter + 1
				// Forward the data to google cloud vertex ai
				forwardDataToGoogleCloudVertexAI(TestJSONValue.Num0.Host, TestJSONValue.Num0.GoogleCloudVertexAiVision.ProjectName, TestJSONValue.Num0.GoogleCloudVertexAiVision.GcpRegion, TestJSONValue.Num0.GoogleCloudVertexAiVision.VertexStreams)
				// Remove a 1 from the counter; the rtsp server has stopped
				rtspConnectionCounter = rtspConnectionCounter - 1
			}
		}
	}
}
