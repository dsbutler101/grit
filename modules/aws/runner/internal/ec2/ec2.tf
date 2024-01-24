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
        content     = var.fleeting.ssh_key_pem
      },
      {
        path        = "/etc/gitlab-runner/config.toml"
        owner       = "root:root"
        permissions = "0755"
        content = templatefile("${path.module}/config.toml", {
          gitlab_url   = var.gitlab.url
          runner_token = var.gitlab.runner_token
          aws_asg_name = var.fleeting.autoscaling_group_name
          executor     = var.executor
          idle_count   = var.scale_min * var.capacity_per_instance
          scale_max    = var.scale_max
        })
      }
    ]

    runcmd = concat(local.install_runner_cmd, var.executor == "docker-autoscaler" || var.executor == "instance" ? local.install_fleeting_plugin_cmd : [])
  }

  install_runner_cmd = [
    "curl -L \"https://packages.gitlab.com/install/repositories/runner/gitlab-runner/script.deb.sh\" | sudo bash",
    "sudo apt-get install gitlab-runner",
  ]

  install_fleeting_plugin_cmd = [
    "curl \"https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip\" -o \"awscliv2.zip\"",
    "apt install unzip",
    "unzip awscliv2.zip",
    "sudo ./aws/install",
    "aws --profile default configure set aws_access_key_id \"${var.iam.fleeting_access_key_id}\"",
    "aws --profile default configure set aws_secret_access_key \"${var.iam.fleeting_secret_access_key}\"",
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

resource "aws_security_group" "manager_sg" {
  name   = "${var.name} manager"
  vpc_id = var.vpc.id

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

  tags = merge(var.labels, {
    Name = "${var.name} manager"
  })
}

resource "aws_instance" "runner-manager" {
  ami                         = data.aws_ami.ubuntu.id
  instance_type               = "t2.micro"
  subnet_id                   = var.vpc.subnet_id
  vpc_security_group_ids      = [aws_security_group.manager_sg.id]
  associate_public_ip_address = true
  user_data                   = data.cloudinit_config.config.rendered

  tags = merge(var.labels, {
    Name = "${var.name}_runner-manager"
  })

  key_name = var.fleeting.ssh_key_pem_name
}

output "public_ip" {
  value = aws_instance.runner-manager.public_ip
}
