# TFLint Ruleset Cinder
[![Build Status](https://github.com/cinder-technologies/tflint-ruleset-cinder/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/cinder-technologies/tflint-ruleset-cinder/actions)

A TFLint plugin that provides additional rules for enforcing Terraform best practices. This ruleset focuses on improving resource naming conventions and code quality.

## Requirements

- TFLint v0.42+
- Go v1.24

## Installation

You can install the plugin with `tflint --init`. Declare a config in `.tflint.hcl` as follows:

```hcl
plugin "cinder" {
  enabled = true

  version = "0.1.1"
  source  = "github.com/cinder-technologies/tflint-ruleset-cinder"
}
```

## Rules

|Name|Description|Severity|Enabled|Link|
| --- | --- | --- | --- | --- |
|terraform_resource_name_contains_type|Prevents resource names from containing redundant type information (shared suffixes with resource type)|ERROR|✔||

### terraform_resource_name_contains_type

This rule enforces cleaner resource naming by preventing resource names from having shared suffixes with their resource types. This helps avoid redundant naming like `aws_s3_bucket "my_bucket"` or `aws_iam_role "user_role"`.

#### Examples

❌ **Bad** - Resource names that share suffixes with their types:
```hcl
resource "aws_s3_bucket" "data_bucket" {
  bucket = "my-data-bucket"
}

resource "aws_iam_role" "user_role" {
  name = "MyUserRole"
}

resource "aws_lambda_function" "my_function" {
  filename = "lambda.zip"
}
```

✅ **Good** - Clean resource names without type redundancy:
```hcl
resource "aws_s3_bucket" "data_storage" {
  bucket = "my-data-bucket"
}

resource "aws_iam_role" "user_permissions" {
  name = "MyUserRole"
}

resource "aws_lambda_function" "data_processor" {
  filename = "lambda.zip"
}
```

The rule performs case-insensitive matching and splits resource types and names by underscores to detect shared suffixes.

## Building the plugin

Clone the repository locally and run the following command:

```
$ make
```

You can easily install the built plugin with the following:

```
$ make install
```

You can run the built plugin like the following:

```
$ cat << EOS > .tflint.hcl
plugin "cinder" {
  enabled = true
}
EOS
$ tflint
```
