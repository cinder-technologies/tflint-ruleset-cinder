package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformResourceNameContainsType(t *testing.T) {
	tests := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "aws_iam_role with full type in name",
			Content: `
resource "aws_iam_role" "user_iam_role" {
  name = "test"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformResourceNameContainsTypeRule(),
					Message: "Resource name 'user_iam_role' contains parts of the resource type 'aws_iam_role'",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 40},
					},
				},
			},
		},
		{
			Name: "aws_iam_role with role in name",
			Content: `
resource "aws_iam_role" "user_role" {
  name = "test"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformResourceNameContainsTypeRule(),
					Message: "Resource name 'user_role' contains parts of the resource type 'aws_iam_role'",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 36},
					},
				},
			},
		},
		{
			Name: "aws_s3_bucket with bucket in name",
			Content: `
resource "aws_s3_bucket" "my_bucket" {
  bucket = "test"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformResourceNameContainsTypeRule(),
					Message: "Resource name 'my_bucket' contains parts of the resource type 'aws_s3_bucket'",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 37},
					},
				},
			},
		},
		{
			Name: "aws_lambda_function with function in name",
			Content: `
resource "aws_lambda_function" "my_function" {
  filename = "test.zip"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformResourceNameContainsTypeRule(),
					Message: "Resource name 'my_function' contains parts of the resource type 'aws_lambda_function'",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 45},
					},
				},
			},
		},
		{
			Name: "case insensitive matching",
			Content: `
resource "aws_iam_role" "user_IAM_ROLE" {
  name = "test"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformResourceNameContainsTypeRule(),
					Message: "Resource name 'user_IAM_ROLE' contains parts of the resource type 'aws_iam_role'",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 40},
					},
				},
			},
		},
		{
			Name: "valid resource name without type",
			Content: `
resource "aws_iam_role" "user_permissions" {
  name = "test"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "valid s3 bucket name",
			Content: `
resource "aws_s3_bucket" "data_storage" {
  bucket = "test"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "valid lambda function name",
			Content: `
resource "aws_lambda_function" "data_processor" {
  filename = "test.zip"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "short provider prefix allowed",
			Content: `
resource "aws_iam_role" "aws_user" {
  name = "test"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "ec2 instance with instance in name",
			Content: `
resource "aws_instance" "web_instance" {
  ami = "ami-12345"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformResourceNameContainsTypeRule(),
					Message: "Resource name 'web_instance' contains parts of the resource type 'aws_instance'",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 39},
					},
				},
			},
		},
		{
			Name: "multiple resources mixed",
			Content: `
resource "aws_iam_role" "user_permissions" {
  name = "test"
}

resource "aws_s3_bucket" "data_bucket" {
  bucket = "test"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformResourceNameContainsTypeRule(),
					Message: "Resource name 'data_bucket' contains parts of the resource type 'aws_s3_bucket'",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 6, Column: 1},
						End:      hcl.Pos{Line: 6, Column: 39},
					},
				},
			},
		},
	}

	rule := NewTerraformResourceNameContainsTypeRule()

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"resource.tf": test.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, test.Expected, runner.Issues)
		})
	}
}