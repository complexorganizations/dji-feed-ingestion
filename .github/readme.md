# DJI Feed Ingestion

Welcome to the DJI Feed Ingestion project by Complex Organizations. This project focuses on ingesting and analyzing data from DJI drones using cloud platforms like AWS, Google Cloud, and Azure.

## 🚀 Join Us!

We're actively hiring designers and developers. If you're passionate about drones and cloud technologies, [reach out](https://complexorganizations.com) to join our dynamic team.

## Requirements

### Hardware

- **Basic**: [DJI Air 3](https://store.dji.com/product/dji-air-3)
- **Professional**: [DJI Mavic 3](https://www.dji.com/mavic-3)
- **Enterprise**: [DJI Matrice 30](https://www.dji.com/matrice-30)

### Software

- AWS, Google Cloud, and Azure

## Setup

1. Establish an RTMP server and deploy cloud resources using Terraform.
2. Use RTMP to relay DJI drone data to cloud service providers.
3. Analyze incoming data with the cloud provider.

## Features

- Accept data from various DJI drones using RTSP.

## Q&A

- **Supported Drones**: Only drones that can be fully automated using FlightHub are supported.
- **Streaming**: Use the DJI App to stream video feed to the cloud.
- **Waypoints Transfer**: Transfer waypoints from DJI flight hub to the controller.
- **Live Stream**: Watch the stream live via VLC.
- **Supported Controllers**: RC-N1, RC-N2, RC-Pro
- **Recommended WISP**: Verizon 5G Home ($60/Month) or T-Mobile 5G Home ($50/Month).
- **Unsupported Controllers**: DJI RC & RC 2 with built-in screen.
- **Platform Preference**: Android (bitrate of 5) over iOS (bitrate of 3).
- **DJI Settings**: Record in 4k 60fps, disable cache for videos, enable subtitles for post-analysis.

## Cloud Services

- **Amazon**: VPC, ELB, EC2, Auto Scaling, Kinesis Video Streams, S3
- **Google Cloud**: VPC, Load Balancing, VMs, Vertex AI Vision, Storage
- **Microsoft Azure**: Virtual Network, Load Balancer, VMs, Blob Storage

## Notes

Ensure a stable internet connection with at least 30 MBPS upload & download speed when connecting to the RTMP server.

## Development

You can develop the code without cloning the repository. However, if you wish to debug, clone the repo and start debugging.

```bash
git clone https://github.com/complexorganizations/dji-feed-ingestion
cd dji-feed-ingestion/
```

## 🤝 Sponsors

[![AWS](https://raw.githubusercontent.com/complexorganizations/parking-united-com/main/assets/images/icons/cloud_providers/aws.svg)](https://aws.amazon.com/)
[![GCP](https://raw.githubusercontent.com/complexorganizations/parking-united-com/main/assets/images/icons/cloud_providers/gcp.svg)](https://cloud.google.com/)
[![Azure](https://raw.githubusercontent.com/complexorganizations/parking-united-com/main/assets/images/icons/cloud_providers/azure.svg)](https://azure.microsoft.com/)
[![DigitalOcean](https://raw.githubusercontent.com/complexorganizations/parking-united-com/main/assets/images/icons/cloud_providers/digitalocean.svg)](https://www.digitalocean.com/)
[![Linode](https://raw.githubusercontent.com/complexorganizations/parking-united-com/main/assets/images/icons/cloud_providers/linode.svg)](https://www.linode.com/)

#### Contact

[![Discord](https://raw.githubusercontent.com/complexorganizations/parking-united-com/main/assets/images/icons/social_media/discord.svg)](https://discord.gg/2DmfdBdMwg)
[![Facebook](https://raw.githubusercontent.com/complexorganizations/parking-united-com/main/assets/images/icons/social_media/facebook.svg)](https://www.facebook.com/parkingunited)
[![Instagram](https://raw.githubusercontent.com/complexorganizations/parking-united-com/main/assets/images/icons/social_media/instagram.svg)](https://www.instagram.com/parking_united/)
[![LinkedIn](https://raw.githubusercontent.com/complexorganizations/parking-united-com/main/assets/images/icons/social_media/linkedin.svg)](https://www.linkedin.com/company/parking-united)
[![Pinterest](https://raw.githubusercontent.com/complexorganizations/parking-united-com/main/assets/images/icons/social_media/pinterest.svg)](https://www.pinterest.com/parking_united/)
[![Reddit](https://raw.githubusercontent.com/complexorganizations/parking-united-com/main/assets/images/icons/social_media/reddit.svg)](https://www.reddit.com/r/parking_united/)
[![Signal](https://raw.githubusercontent.com/complexorganizations/parking-united-com/main/assets/images/icons/social_media/signal.svg)](https://signal.group/#CjQKIPhEy6Pk8c-wXi-6O3DRXQ3eSLvJNqW61uq46Y-Ya3mrEhDaILflpc1oE9joFmzC3REG)
[![Skype](https://raw.githubusercontent.com/complexorganizations/parking-united-com/main/assets/images/icons/social_media/skype.svg)](https://join.skype.com/hjhsrvQlinZk)
[![Slack](https://raw.githubusercontent.com/complexorganizations/parking-united-com/main/assets/images/icons/social_media/slack.svg)](https://parking-unitedcom.slack.com/archives/C05QM7PS9GV/p1693631754500589)
[![Snapchat](https://raw.githubusercontent.com/complexorganizations/parking-united-com/main/assets/images/icons/social_media/snapchat.svg)](https://www.snapchat.com/)
[![Teams](https://raw.githubusercontent.com/complexorganizations/parking-united-com/main/assets/images/icons/social_media/teams.svg)](https://teams.live.com/l/community/FAAHt8haBHMqRRUOwI)
[![Telegram](https://raw.githubusercontent.com/complexorganizations/parking-united-com/main/assets/images/icons/social_media/telegram.svg)](https://t.me/parking_united_com)
[![TikTok](https://raw.githubusercontent.com/complexorganizations/parking-united-com/main/assets/images/icons/social_media/tiktok.svg)](https://www.tiktok.com/)
[![Twitter](https://raw.githubusercontent.com/complexorganizations/parking-united-com/main/assets/images/icons/social_media/twitter.svg)](https://twitter.com/parking_united)
[![WhatsApp](https://raw.githubusercontent.com/complexorganizations/parking-united-com/main/assets/images/icons/social_media/whatsapp.svg)](https://chat.whatsapp.com/KR0nia4ajom2NWl32YOYZK)
[![YouTube](https://raw.githubusercontent.com/complexorganizations/parking-united-com/main/assets/images/icons/social_media/youtube.svg)](https://www.youtube.com/)

## Author

- **Name**: Complex Organizations
- **Website**: [complexorganizations.com](https://complexorganizations.com)

## Credits

A big shoutout to the open-source community for their contributions.

## License

This project is licensed under the [Apache License Version 2.0](https://raw.githubusercontent.com/complexorganizations/dji-feed-ingestion/main/.github/license).
