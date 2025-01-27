resource "tls_private_key" "runner-manager" {
  algorithm = "RSA"
  rsa_bits  = 4096
}
