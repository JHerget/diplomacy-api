terraform {
  required_version = ">= 1.4.0"

  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = ">= 5.0"
    }
    archive = {
      source = "hashicorp/archive"
      version = ">= 2.4"
    }
    local = {
      source = "hashicorp/local"
      version = ">= 2.4"
    }
  }
}

provider "aws" {
  region = "us-west-2"
}
