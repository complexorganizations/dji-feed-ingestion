terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.0"
    }
    google = {
      source  = "hashicorp/google"
      version = ">= 4.43.1"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 3.33.0"
    }
  }
  required_version = ">= 1.0.0"
}

// Configure the AWS provider
provider "aws" {
  region                   = "us-east-1"
  shared_credentials_files = ["~/.aws/credentials"]
}

// Configure the GCP provider
provider "google" {
  project = "github-code-snippets"
  region  = "us-central1"
  zone    = "us-central1-a"
}

// Configure the Azure provider
provider "azurerm" {
  features {}
}

// Check which provider is configured properly and only use that provider;
// If multiple providers are configured than use mutiple ones.
