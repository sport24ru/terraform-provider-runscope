# Terraform Runscope Provider

The Runscope provider is used to interact with the resources
supported by [Runscope](https://runscope.com/).

## Installing

### terraform 0.13+

Add into your Terraform configuration this code:

```hcl-terraform
terraform {
  required_providers {
    runscope = {
      source = "sport24ru/runscope"
    }
  }
}
```

and run `terraform init`

### terraform 0.12 and earlier

1. Download archive with the latest version of provider for your operating system from
   [Github releases page](https://github.com/sport24ru/terraform-provider-runscope/releases).
2. Unpack provider to `$HOME/.terraform.d/plugins`, i.e.
   ```
   unzip terraform-provider-runscope_vX.Y.Z_linux_amd64.zip terraform-provider-runscope_* -d $HOME/.terraform.d/plugins/
   ```
3. Init your terraform project
   ```
   terraform init
   ```

## Usage

Read the [documentation on Terraform Registry site](https://registry.terraform.io/providers/sport24ru/runscope/latest/docs).
