resource "google_kms_key_ring" "default" {
  name     = var.metadata.name
  location = var.kms_location
}

resource "google_kms_crypto_key" "default" {
  name            = var.metadata.name
  key_ring        = google_kms_key_ring.default.id
  rotation_period = "604800s"
}

resource "google_kms_secret_ciphertext" "runner_token" {
  crypto_key = google_kms_crypto_key.default.id
  plaintext  = var.runner_token
}

resource "google_kms_secret_ciphertext" "runner_ssh_key" {
  crypto_key = google_kms_crypto_key.default.id
  plaintext  = tls_private_key.runner_manager.private_key_openssh
}
