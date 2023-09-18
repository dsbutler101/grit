locals {
  message = "There is nothing to apply here."
  ok      = false
}
check "ok" {
  assert {
    condition     = local.ok == true
    error_message = local.message
  }
}
