## How to access the terraform state

```shell
export TF_HTTP_USERNAME="<GITLAB USERNAME>"
export TF_HTTP_PASSWORD="<GITLAB PAT WITH API SCOPE>"
terraform init
```

## How to apply

```shell
export GITLAB_TOKEN="<GITLAB PAT WITH API SCOPE>"

# or use something like aws-vault
export AWS_ACCESS_KEY_ID="your user in the shared runner AWS cloud sandbox"
export AWS_SECRET_ACCESS_KEY="..."

# or use gcloud auth
export GOOGLE_APPLICATION_CREDENTIALS="credentials for your user in the shared runner GCP cloud sandbox"

terraform apply
```
