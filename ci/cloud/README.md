## How to access the terraform state

```shell
export TF_HTTP_USERNAME="<GITLAB USERNAME>"
export TF_HTTP_PASSWORD="<GITLAB PAT WITH API SCOPE>"
terraform init
```

## How to apply

```shell
export GITLAB_TOKEN="<GITLAB PAT WITH API SCOPE>"
terraform apply
```
