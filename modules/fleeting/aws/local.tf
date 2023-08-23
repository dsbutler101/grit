locals {
  # AccountID of our Verify Runner sandbox account
  # It is used in maintainer-access.tf to define a role that can be used
  # to authenticate into the SaaS MacOS environment account through the
  # Verify Runner sandbox account
  eng_dev_verify_runner_account_id = 915502504722

  # This is a static CIDR from the subnet in GCP where all our runner
  # managers are hosted.
  # It is used in networking.tf to define firewall rules for accessing
  # ASG instances
  gcp_runner_manager_subnetwork_cidr = "10.1.5.0/24"

  # Set of tags added to most of the resources managed by this module
  tags = {
    gl_realm              = var.realm
    gl_env_type           = var.env_type
    gl_env_name           = var.shard
    gl_owner_email_handle = var.gl_owner_email_handle
    gl_dept               = var.gl_dept
    gl_dept_group         = var.gl_dept_group
    shard                 = var.shard
  }
}
