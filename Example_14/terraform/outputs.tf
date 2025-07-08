output "instance_public_ip" {
  value = aws_instance.back-end_app.public_ip
}

