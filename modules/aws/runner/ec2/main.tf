locals {
  # Note: Do not write files under /tmp during boot 
  # https://cloudinit.readthedocs.io/en/latest/reference/modules.html#write-files
  cloudwatch_agent_json_path = "/var/tmp/amazon-cloudwatch-agent.json"

  cloud_config = {
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
          runner_name             = var.name
          aws_asg_name            = var.fleeting.autoscaling_group_name
          username                = var.fleeting.username
          executor                = var.executor
          runners_global_section  = var.runners_global_section
          idle_count              = var.scale_min * var.capacity_per_instance
          idle_time               = var.idle_time
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
          runners_docker_section  = var.runners_docker_section
          usage_logger            = var.usage_logger
          s3_cache                = var.s3_cache
          acceptable_durations    = var.acceptable_durations
        })
      },
      ],
      var.runner_wrapper.enabled ? [{
        path        = "/etc/systemd/system/gitlab-runner.service.d/wrapper.conf"
        owner       = "root:root"
        permissions = "0644"
        content     = templatefile("${path.module}/wrapper.conf", var.runner_wrapper)
      }] : [],
      var.node_exporter.enabled ? module.node_exporter[0].write_files_config : [],
      var.install_cloudwatch_agent ? [local.cloudwatch_config_file] : []
    )

    runcmd = concat(
      var.install_cloudwatch_agent ? local.install_cloudwatch_agent_cmd : [],
      var.node_exporter.enabled ? module.node_exporter[0].commands : [],
      local.install_runner_cmd,
      var.executor == "docker-autoscaler" || var.executor == "instance" ? local.install_fleeting_plugin_cmd : [],
    )
  }

  cloudwatch_config_file = {
    path        = local.cloudwatch_agent_json_path
    owner       = "root:root"
    permissions = "0644"
    content     = base64decode(var.cloudwatch_agent_json)
  }

  install_cloudwatch_agent_cmd = [
    "sudo curl -O https://amazoncloudwatch-agent.s3.amazonaws.com/ubuntu/amd64/latest/amazon-cloudwatch-agent.deb",
    "sudo apt-get -o DPkg::Lock::Timeout=300 install ./amazon-cloudwatch-agent.deb -y",
    "sudo usermod -aG adm cwagent",
    "sudo amazon-cloudwatch-agent-ctl -a start",
    "sudo amazon-cloudwatch-agent-ctl -a fetch-config -c file:${local.cloudwatch_agent_json_path} -s"
  ]

  install_runner_cmd = [
    "export PATH=\"/etc/gitlab-runner:$PATH\"",
    "curl -L \"https://packages.gitlab.com/install/repositories/runner/${var.runner_repository}/script.deb.sh\" | sudo bash",
    "sudo apt-get install -ym gitlab-runner-helper-images=${var.runner_version}",
    "sudo apt-get install -y gitlab-runner=${var.runner_version}",
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

  instance_name = "${var.name}_runner-manager"
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

module "node_exporter" {
  count                 = var.node_exporter.enabled ? 1 : 0
  source                = "../../node_exporter"
  node_exporter_port    = var.node_exporter.port
  node_exporter_version = var.node_exporter.version
}

resource "aws_instance" "runner_manager" {
  ami                         = var.runner_manager_ami != "" ? var.runner_manager_ami : data.aws_ami.ubuntu.id
  instance_type               = var.instance_type
  subnet_id                   = try(length(var.vpc.subnet_ids), 0) > 0 ? var.vpc.subnet_ids[0] : var.vpc.subnet_id
  vpc_security_group_ids      = var.security_group_ids
  associate_public_ip_address = var.associate_public_ip_address
  user_data                   = data.cloudinit_config.config.rendered
  iam_instance_profile        = var.instance_role_profile_name
  user_data_replace_on_change = true
  metadata_options {
    http_endpoint = "enabled"
    http_tokens   = "required"
  }

  tags = merge(var.labels, {
    Name = local.instance_name
  })

  key_name = try(aws_key_pair.aws_runner_key_pair[0].key_name, "")

  root_block_device {
    delete_on_termination = true
    throughput            = var.throughput
    volume_size           = var.volume_size
    volume_type           = var.volume_type
    encrypted             = var.encrypted
    kms_key_id            = var.kms_key_id
  }
}
