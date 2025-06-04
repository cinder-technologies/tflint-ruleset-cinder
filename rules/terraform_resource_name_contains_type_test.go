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
			Name: "aws_iam_role with shared suffix iam_role",
			Content: `
resource "aws_iam_role" "my_iam_role" {
  name = "test"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformResourceNameContainsTypeRule(),
					Message: "Resource name 'my_iam_role' should not have a shared suffix with the resource type 'aws_iam_role'.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 38},
					},
				},
			},
		},
		{
			Name: "aws_iam_role with shared suffix role",
			Content: `
resource "aws_iam_role" "user_role" {
  name = "test"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformResourceNameContainsTypeRule(),
					Message: "Resource name 'user_role' should not have a shared suffix with the resource type 'aws_iam_role'.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 36},
					},
				},
			},
		},
		{
			Name: "aws_iam_role with no shared suffix - allowed",
			Content: `
resource "aws_iam_role" "my_iam" {
  name = "test"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "aws_s3_bucket with shared suffix bucket",
			Content: `
resource "aws_s3_bucket" "my_bucket" {
  bucket = "test"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformResourceNameContainsTypeRule(),
					Message: "Resource name 'my_bucket' should not have a shared suffix with the resource type 'aws_s3_bucket'.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 37},
					},
				},
			},
		},
		{
			Name: "aws_lambda_function with shared suffix function",
			Content: `
resource "aws_lambda_function" "my_function" {
  filename = "test.zip"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformResourceNameContainsTypeRule(),
					Message: "Resource name 'my_function' should not have a shared suffix with the resource type 'aws_lambda_function'.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 45},
					},
				},
			},
		},
		{
			Name: "case insensitive matching - shared suffix iam_role",
			Content: `
resource "aws_iam_role" "user_IAM_ROLE" {
  name = "test"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformResourceNameContainsTypeRule(),
					Message: "Resource name 'user_IAM_ROLE' should not have a shared suffix with the resource type 'aws_iam_role'.",
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
			Name: "resource name with shared suffix policy - should not trigger",
			Content: `
resource "aws_iam_role_policy_attachment" "cloudwatch_agent_server_policy" {
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
			Name: "aws_instance with shared suffix instance",
			Content: `
resource "aws_instance" "web_instance" {
  ami = "ami-12345"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformResourceNameContainsTypeRule(),
					Message: "Resource name 'web_instance' should not have a shared suffix with the resource type 'aws_instance'.",
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
					Message: "Resource name 'data_bucket' should not have a shared suffix with the resource type 'aws_s3_bucket'.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 6, Column: 1},
						End:      hcl.Pos{Line: 6, Column: 39},
					},
				},
			},
		},
		{
			Name: "type part in middle of name should be allowed",
			Content: `
resource "aws_iam_role" "role_for_lambda" {
  name = "test"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "type part in middle of name should be allowed - bucket",
			Content: `
resource "aws_s3_bucket" "bucket_for_data" {
  bucket = "test"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "resource name with shared suffix policy - should not trigger",
			Content: `
resource "aws_iam_role_policy_attachment" "amazon_eks_cluster_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "resource name with shared suffix attachment - should trigger",
			Content: `
resource "aws_iam_role_policy_attachment" "amazon_eks_cluster_policy_attachment" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformResourceNameContainsTypeRule(),
					Message: "Resource name 'amazon_eks_cluster_policy_attachment' should not have a shared suffix with the resource type 'aws_iam_role_policy_attachment'.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 81},
					},
				},
			},
		},
		{
			Name: "resource name with no shared suffix - allowed",
			Content: `
resource "aws_iam_role_policy_attachment" "eks_cluster" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "aws_security_group with shared suffix group",
			Content: `
resource "aws_security_group" "web_security_group" {
  name = "web-sg"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformResourceNameContainsTypeRule(),
					Message: "Resource name 'web_security_group' should not have a shared suffix with the resource type 'aws_security_group'.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 51},
					},
				},
			},
		},
		{
			Name: "nonshared suffix like s3 in name is allowed",
			Content: `
resource "aws_s3_bucket" "logs_s3" {
  bucket = "my-logs"
}`,
			Expected: helper.Issues{},
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