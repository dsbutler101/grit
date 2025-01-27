resource "google_kms_key_ring" "default" {
  name     = var.metadata.name
  location = "global"
}

resource "google_kms_crypto_key" "default" {
  name            = var.metadata.name
  key_ring        = google_kms_key_ring.default.id
  rotation_period = "604800s"
}

resource "google_kms_secret_ciphertext" "runner-token" {
  crypto_key = google_kms_crypto_key.default.id
  plaintext  = var.runner_token
}

resource "google_kms_secret_ciphertext" "runner-ssh-key" {
  crypto_key = google_kms_crypto_key.default.id
  plaintext  = tls_private_key.runner-manager.private_key_openssh
}
