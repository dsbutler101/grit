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

    write_files = setunion([
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
          concurrent              = var.scale_max * var.capacity_per_instance
          gitlab_url              = var.gitlab.url
          runner_token            = var.gitlab.runner_token
          aws_asg_name            = var.fleeting.autoscaling_group_name
          username                = var.fleeting.username
          executor                = var.executor
          idle_count              = var.scale_min * var.capacity_per_instance
          scale_max               = var.scale_max
          max_use_count           = var.max_use_count
          idle_percentage         = var.idle_percentage
          privileged              = var.privileged
          region                  = var.region
          enable_metrics_export   = var.enable_metrics_export
          metrics_export_endpoint = var.metrics_export_endpoint
          aws_plugin_version      = var.aws_plugin_version
          capacity_per_instance   = var.capacity_per_instance
          default_docker_image    = var.default_docker_image
          s3_cache                = var.s3_cache
        })
      }
      ],
      var.install_cloudwatch_agent ? [local.cloudwatch_config_file] : []
    )


    runcmd = concat(
      var.install_cloudwatch_agent ? local.install_cloudwatch_agent_cmd : [],
      local.install_runner_cmd,
      var.executor == "docker-autoscaler" || var.executor == "instance" ? local.install_fleeting_plugin_cmd : [],
    )
  }

  cloudwatch_config_file = {
    path        = "/tmp/amazon-cloudwatch-agent.json"
    owner       = "root:root"
    permissions = "0644"
    content     = base64decode(var.cloudwatch_agent_json)
  }

  install_cloudwatch_agent_cmd = [
    "sudo curl -O https://amazoncloudwatch-agent.s3.amazonaws.com/ubuntu/amd64/latest/amazon-cloudwatch-agent.deb",
    "sudo apt-get -o DPkg::Lock::Timeout=300 install ./amazon-cloudwatch-agent.deb -y",
    "sudo usermod -aG adm cwagent",
    "sudo amazon-cloudwatch-agent-ctl -a start",
    "sudo amazon-cloudwatch-agent-ctl -a fetch-config -c file:/tmp/amazon-cloudwatch-agent.json -s"
  ]

  install_runner_cmd = [
    "export PATH=\"/etc/gitlab-runner:$PATH\"",
    "curl -L \"https://packages.gitlab.com/install/repositories/runner/${var.runner_repository}/script.deb.sh\" | sudo bash",
    "sudo apt-get install gitlab-runner=${var.runner_version}",
  ]

  install_fleeting_plugin_cmd = [
    "curl \"https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip\" -o \"awscliv2.zip\"",
    "apt install unzip",
    "unzip awscliv2.zip",
    "sudo ./aws/install",
    "aws --profile default configure set aws_access_key_id \"${var.iam.fleeting_access_key_id}\"",
    "aws --profile default configure set aws_secret_access_key \"${var.iam.fleeting_secret_access_key}\"",
    "aws --profile default configure set region \"${var.region}\"",
    "gitlab-runner fleeting install",
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

resource "aws_instance" "runner-manager" {
  ami                         = data.aws_ami.ubuntu.id
  instance_type               = "t2.micro"
  subnet_id                   = var.vpc.subnet_id
  vpc_security_group_ids      = var.security_group_ids
  associate_public_ip_address = true
  user_data                   = data.cloudinit_config.config.rendered
  iam_instance_profile        = var.instance_role_profile_name
  user_data_replace_on_change = true

  tags = merge(var.labels, {
    Name = "${var.name}_runner-manager"
  })

  key_name = var.fleeting.ssh_key_pem_name
}

output "public_ip" {
  value = aws_instance.runner-manager.public_ip
}
