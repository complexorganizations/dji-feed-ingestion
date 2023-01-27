terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.0"
    }
  }
  required_version = ">= 1.0.0"
}

// Configure the AWS Provider
provider "aws" {
  region                   = "us-east-1"
  shared_credentials_files = ["~/.aws/credentials"]
}

// Create a Kinesis video stream
resource "aws_kinesis_video_stream" "kinesis_video_stream" {
  // {project-name}-kinesis-video-{0}-{us-east-1}
  name                    = "code-snippets-kinesis-video-0-us-east-1"
  data_retention_in_hours = 7
  media_type              = "video/h264"
  tags = {
    Name = "code-snippets-kinesis-video-0-us-east-1"
  }
}
