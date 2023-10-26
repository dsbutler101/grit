locals {
  cloud_config = {
    groups = [
      {
        ubuntu = ["root", "sys"]
      },
      "hashicorp"
    ]

    users = ["default"]

    packages = ["git", "git-lfs"]

    write_files = [
      {
        path        = "/etc/gitlab-runner/keypair.pem"
        owner       = "root:root"
        permissions = "400"
        content     = var.ssh_key_pem
      },
      {
        path        = "/etc/gitlab-runner/config.toml"
        owner       = "root:root"
        permissions = "0755"
        content = templatefile("${path.module}/config.toml", {
          gitlab_url        = var.gitlab_url
          runner_token      = var.runner_token
          aws_asg_name      = var.aws_asg_name
          scale_max         = var.scale_max
          idle_count        = var.idle_count
          fleeting_provider = var.fleeting_provider
        })
      }
    ]

    runcmd = concat(local.basic_cmd, var.fleeting_provider == "ec2" ? local.aws_fleet_cmd : [])
  }

  basic_cmd = [
    "curl -L \"https://packages.gitlab.com/install/repositories/runner/gitlab-runner/script.deb.sh\" | sudo bash",
    "sudo apt-get install gitlab-runner"
  ]

  aws_fleet_cmd = [
    "curl \"https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip\" -o \"awscliv2.zip\"",
    "apt install unzip",
    "unzip awscliv2.zip",
    "sudo ./aws/install",
    "aws --profile default configure set aws_access_key_id \"${var.fleeting_service_account_access_key_id}\"",
    "aws --profile default configure set aws_secret_access_key \"${var.fleeting_service_account_secret_access_key}\"",
    "aws --profile default configure set region \"us-east-1\"",
    "curl -Lo /etc/gitlab-runner/fleeting-plugin-aws \"https://gitlab.com/gitlab-org/fleeting/fleeting-plugin-aws/-/releases/permalink/latest/downloads/fleeting-plugin-aws-linux-amd64\"",
    "chmod +x /etc/gitlab-runner/fleeting-plugin-aws && chown gitlab-runner /etc/gitlab-runner/fleeting-plugin-aws",
    "chown gitlab-runner /etc/gitlab-runner/keypair.pem"
  ]

  rendered_yaml = yamlencode(local.cloud_config)
}

data "cloudinit_config" "config" {
  gzip          = false
  base64_encode = false

  part {
    filename     = "cloud-config.yaml"
    content_type = "text/cloud-config"

    content = local.rendered_yaml
  }
}

data "aws_ami" "ubuntu" {
  most_recent = true
  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-*20*-amd64-server-*"]
  }
  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
  owners = ["099720109477"] # Canonical
}

resource "aws_vpc" "vpc" {
  cidr_block           = "10.1.0.0/16"
  enable_dns_support   = true
  enable_dns_hostnames = true
}

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.vpc.id
}

resource "aws_subnet" "subnet_public" {
  vpc_id     = aws_vpc.vpc.id
  cidr_block = "10.1.0.0/24"
}

resource "aws_route_table" "rtb_public" {
  vpc_id = aws_vpc.vpc.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
  }
}

resource "aws_route_table_association" "rta_subnet_public" {
  subnet_id      = aws_subnet.subnet_public.id
  route_table_id = aws_route_table.rtb_public.id
}

resource "aws_security_group" "sg_22" {
  name   = "sg_22"
  vpc_id = aws_vpc.vpc.id
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "runner-manager" {
  ami                         = data.aws_ami.ubuntu.id
  instance_type               = "t2.micro"
  subnet_id                   = aws_subnet.subnet_public.id
  vpc_security_group_ids      = [aws_security_group.sg_22.id]
  associate_public_ip_address = true
  user_data                   = data.cloudinit_config.config.rendered
  tags = {
    Name = "GRIT Runner Manager"
  }
  key_name = var.ssh_key_pem_name
}

output "public_ip" {
  value = aws_instance.runner-manager.public_ip
}
