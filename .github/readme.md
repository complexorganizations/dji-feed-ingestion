### DJI Feed Ingestion

---
### Contact:

[![Discord](https://raw.githubusercontent.com/complexorganizations/parking-united-com/main/assets/discord.svg)](https://discord.gg/3mwJHr5Btt)
[![Slack](https://raw.githubusercontent.com/complexorganizations/parking-united-com/main/assets/slack.svg)](https://parking-unitedcom.slack.com/archives/C05QM7PS9GV/p1693631754500589)
[![Signal](https://raw.githubusercontent.com/complexorganizations/parking-united-com/main/assets/signal.svg)](https://signal.group/#CjQKIPhEy6Pk8c-wXi-6O3DRXQ3eSLvJNqW61uq46Y-Ya3mrEhDaILflpc1oE9joFmzC3REG)
[![WhatsApp](https://raw.githubusercontent.com/complexorganizations/parking-united-com/main/assets/whatsapp.svg)](https://chat.whatsapp.com/KR0nia4ajom2NWl32YOYZK)
[![Telegram](https://raw.githubusercontent.com/complexorganizations/parking-united-com/main/assets/telegram.svg)](https://t.me/parking_united_com)
[![Teams](https://raw.githubusercontent.com/complexorganizations/parking-united-com/main/assets/teams.svg)](https://teams.live.com/l/community/FAAHt8haBHMqRRUOwI)
[![Skype](https://raw.githubusercontent.com/complexorganizations/parking-united-com/main/assets/skype.svg)](https://join.skype.com/hjhsrvQlinZk)
[![Reddit](https://raw.githubusercontent.com/complexorganizations/parking-united-com/main/assets/reddit.svg)](https://reddit.com/r/parking_united_com/s/YefTKn0cfX)


---
### Requirements:

Hardware requirements:

- Professional: DJI Mavic 3 `https://www.dji.com/mavic-3`
- Enterprise: DJI Matrice 30 & Dock `https://www.dji.com/matrice-30`

Software requirements:

- AWS & Google Cloud & Azure

---
### Overall Setup:
- Establish a RTMP server and deploy cloud resources using terraform.
- Use RTMP to send all of the DJI drone's data to the cloud service providers.
- Use the cloud provider to analyse all incoming data.

---
### Features
- Accept data from various DJI drones using RTSP.
- 
 
---
### Q&A

#### Why are only these drone supported?
- Currently only these drones are supported since the flight path can be fully automated using [FlightHub](https://www.dji.com/flighthub-2)

#### How to stream the video feed from your DJI drone to the cloud?
- DJI APP > Settings > Transmission > Live Streaming Platforms > RTMP > `rtmp://localhost:1935/drone_zero?user=PublishAdministrator&pass=PublishPassword`

#### How to transfer waypoints from DJI flight hub to controller?
- ``

#### How to watch the stream live via vlc?
- VLC APP > Media > Open Network Stream > `rtsp://ReadAdministrator:ReadPassword@localhost:8554/test_zero`

#### Which controlls are supported?
- RC-N1
- DJI RC Pro

#### Which WISP would you recommenend?
- [Verizon 5G Home](https://www.verizon.com/5g/home/) $60/Month
- [T-Mobile 5G Home](https://www.t-mobile.com/home-internet) $50/Month

#### Why are the other conrollers not supported?
- DJI RC ***The one with the screen build in is NOT supported.***

#### Is it better to use android or ios?
- Android has a bitrate of 5 while ios has a bitrate of 3. ***Android***

#### What settings should the DJI use?
- Disable cache for videos *Useless storage*
- Enable subtitles for videos (location; gps cordnates)... post analysis.
- Record in 4k 60fps; use the auto feature.

#### Which cloud services are used?

##### [Amazon](https://aws.amazon.com/)
- Amazon Virtual Private Cloud
- Amazon Elastic Load Balancer (ELB)
- Amazon Elastic Compute Cloud
- Amazon Auto Scaling
- Amazon Kinesis Video Streams
- Amazon S3

##### [Google Cloud](https://cloud.google.com/)
- Virtual Private Cloud (VPC) | Google Cloud
- Cloud Load Balancing | Google Cloud
- Compute Engine Virtual Machines (VMs) | Google Cloud
- Vertex AI Vision | Google Cloud
- Cloud Storage | Google Cloud

##### [Microsoft Azure](https://azure.microsoft.com/)
- Microsoft Azure Virtual Network
- Microsoft Azure Load Balancer
- Microsoft Azure Virtual Machines
- Microsoft Azure "REPLACE_THIS_WITH_KVS_OR_IVS_OR_VERTEX_AI_ON_AZURE"
- Microsoft Azure Blob Storage

---
### Notes
- Make sure you have 30 MBPS Upload & 30 MBPS Download when connecting to the RTMP server.

---
### Author
* Name: Complex Organizations
* Website: [complexorganizations.com](https://www.complexorganizations.com)

---	
### Credits
- Open Source Community

---
### License
- Apache License Version 2.0
