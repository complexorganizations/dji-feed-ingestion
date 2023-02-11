sudo apt-get install apt-transport-https ca-certificates gnupg -y
echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | sudo tee -a /etc/apt/sources.list.d/google-cloud-sdk.list
curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key --keyring /usr/share/keyrings/cloud.google.gpg add -
sudo apt-get update && sudo apt-get install google-cloud-cli -y
gcloud init
gcloud services enable visionai.googleapis.com
wget https://github.com/google/visionai/releases/download/v0.0.4/visionai_0.0-4_amd64.deb
sudo apt install ./visionai_0.0-4_amd64.deb
sudo apt install bazel-4.2.1 -y
apt-get install -y --no-install-recommends autoconf automake build-essential ca-certificates flex bison python3 nasm libjpeg-dev -y
# https://cloud.google.com/vision-ai/docs/create-manage-streams#ingest-videos
# This command will send an RTSP feed into the stream.
# This command has to run in the network that has direct access to the RTSP feed.
vaictl -p PROJECT_ID \
         -l LOCATION_ID \
         -c application-cluster-0 \
         --service-endpoint visionai.googleapis.com \
send rtsp to streams STREAM_ID --rtsp-uri RTSP_ADDRESS
