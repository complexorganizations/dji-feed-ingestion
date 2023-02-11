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

// Create 3 kinesis video stream
// Create 1 s3 bucket
// Create 1 DynamoDB
