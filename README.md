# Diplomacy API

## Terraform state

Terraform is configured to store state in S3 at
`s3://diplomacy-api-terraform-state/diplomacy-api/terraform.tfstate`.

Create the state bucket before initializing Terraform:

```sh
aws s3api create-bucket \
  --bucket diplomacy-api-terraform-state \
  --region us-west-2 \
  --create-bucket-configuration LocationConstraint=us-west-2
```

Then migrate the current local state:

```sh
cd terraform
terraform init -migrate-state
```
