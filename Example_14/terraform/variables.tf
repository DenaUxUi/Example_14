variable "aws_region" {
	default = "eu-north-1"
}
variable "ami_id" {
	description = "Amazon Machine Image ID"
	default	    = "ami-042b4708b1d05f512"
}
variable "instance_type" {
	default = "t2.micro"
}
variable "key_name" {
	description = "Name of SSH-key from EC2 > Key Pairs"
	default = "Dena"
} 
