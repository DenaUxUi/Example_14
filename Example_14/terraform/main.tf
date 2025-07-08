terraform {

  required_providers {

    aws = {

      source  = "hashicorp/aws"

      version = "~> 6.0"

    }

  }

  required_version = ">= 1.2.0"

}
provider "aws" {

  region  = "eu-north-1"

  profile = "default"

}



resource "aws_instance" "back-end_app" {

  ami           = "ami-042b4708b1d05f512"

  instance_type = "t3.micro"
  key_name      = "Dena"


  tags = {

    Name = "DenaInstance"

  }

}
