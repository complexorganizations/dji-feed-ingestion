# DJI Data Integration Initiative

Embark on the DJI Data Integration Initiative, presented by Complex Organizations. Our mission centers on the ingestion and interpretation of data from DJI drones, utilizing prominent cloud platforms such as AWS, Google Cloud, and Azure to achieve our goals.

## 🚀 Careers at Complex Organizations

Join our expanding team! We are looking for skilled designers and developers with a zeal for drone technology and cloud solutions. If that sounds like you, don’t hesitate to [get in touch](https://complexorganizations.com) and become part of our vibrant team.

## Project Prerequisites

### Equipment Specifications

- **Entry-Level**: [DJI Air 3](https://www.dji.com/air-3)
- **Advanced-Level**: [DJI Mavic 3 Classic](https://www.dji.com/mavic-3-classic)
- **Industrial-Level**: [DJI Dock 2](https://enterprise.dji.com/dock-2)

### Software and Platforms

- Involves the use of AWS, Google Cloud, and Azure ecosystems

## Implementation Strategy

1. Initiate by setting up an RTMP server and orchestrating cloud resources through Terraform.
2. Employ RTMP for the seamless transmission of data from DJI drones to our cloud partners.
3. Engage in comprehensive data analysis via the cloud platform's tools.

## Key Functionalities

- Capability to receive and process data from a range of DJI drones using RTSP protocol.

## Frequently Asked Questions (FAQ)

### Drone Compatibility
- **Eligible Drones**: Compatibility is limited to drones that offer full automation via FlightHub.
- **Streaming Functionality**: Utilize the DJI App for direct video feed streaming to cloud platforms.
- **Waypoint Synchronization**: Seamlessly transfer waypoints from DJI Flight Hub to the drone controller.
- **Real-Time Streaming**: Experience live streaming capabilities through VLC player.
- **Controller Support**: Compatible with RC-N1, RC-N2, and RC-Pro controllers.
- **Wireless Service Recommendations**: Suggested WISPs are Verizon 5G Home (approx. $60/Month) and T-Mobile 5G Home (around $50/Month).
- **Non-Supported Controllers**: DJI RC and RC 2 controllers with integrated screens are not supported.
- **Preferred Operating System**: Android, due to its higher bitrate efficiency (5) compared to iOS (3).
- **Optimal DJI Drone Settings**: Recommended to record in 4K at 60fps, turn off video caching, and enable subtitles for detailed post-analysis.

### Cloud Integration and Services

#### Amazon Web Services (AWS)
- **Components**: Virtual Private Cloud (VPC), Elastic Load Balancer (ELB), EC2 instances, Auto Scaling capabilities, Kinesis Video Streams, and S3 storage solutions.

#### Google Cloud Platform
- **Features**: Virtual Private Cloud (VPC), Load Balancing mechanisms, Virtual Machines (VMs), Vertex AI for Vision capabilities, and extensive Storage options.

#### Microsoft Azure
- **Services**: Implementation of Virtual Networks, Load Balancers, Azure Virtual Machines, and Blob Storage facilities.

## Important Considerations

### Internet Connectivity
- **Network Stability**: For optimal performance, maintain a stable internet connection.
- **Speed Requirement**: Ensure a minimum of 30 MBPS for both upload and download speeds when interfacing with the RTMP server. This is crucial for reliable data transfer and streaming quality.

## Development Guidelines

### Working on the Project
- **Non-Local Development**: Directly contribute to the codebase without needing to clone the repository, allowing for flexibility and immediacy in development.
- **Debugging and In-Depth Development**:
  - For detailed debugging or extensive development work, cloning the repository is recommended. Follow these steps to set up your local development environment:
    ```bash
    git clone https://github.com/complexorganizations/dji-feed-ingestion
    cd dji-feed-ingestion/
    ```
  - This approach facilitates a deeper engagement with the project, allowing developers to test and modify the codebase extensively.

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

## Project Authorship

### Complex Organizations
- **About Us**: We are a team committed to innovative solutions in drone technology and cloud computing.
- **Explore More**: Visit our official website for further insights into our projects and ethos: [complexorganizations.com](https://complexorganizations.com).

## Acknowledgements

### Community Contributions
- **Heartfelt Thanks**: We extend our deepest gratitude to the open-source community. Your contributions play a pivotal role in the success and continual development of our projects.

## Licensing Information

### Open Source License
- **License Type**: The DJI Feed Ingestion project is under the Apache License Version 2.0.
- **Full License Text**: For detailed terms and conditions, please refer to the [Apache License Version 2.0 documentation](https://raw.githubusercontent.com/complexorganizations/dji-feed-ingestion/main/.github/license).
