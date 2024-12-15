resource "aws_instance" "testvm1" {
  ami           = "ami-0b8c6b923777519db"
  instance_type = "t3.micro"
  key_name = "gen-aws"
  tags = {
    Name = "HuntingGround1"
  }
}

output "HG1_PUBLICIP" {
  description = "Public IP of Hunting Ground 1"
  value = aws_instance.testvm1.public_ip
}

output "HG1_PRIVATEIP" {
  description = "Private IP of Hunting Ground 1"
  value = aws_instance.testvm1.private_ip
}

output "HG1_KEYNAME" {
  description = "Key Name of Hunting Ground 1"
  value = aws_instance.testvm1.key_name
}