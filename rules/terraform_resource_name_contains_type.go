package rules

import (
	"fmt"
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformResourceNameContainsTypeRule checks whether resource names contain parts of the resource type
type TerraformResourceNameContainsTypeRule struct {
	tflint.DefaultRule
}

// NewTerraformResourceNameContainsTypeRule returns a new rule
func NewTerraformResourceNameContainsTypeRule() *TerraformResourceNameContainsTypeRule {
	return &TerraformResourceNameContainsTypeRule{}
}

// Name returns the rule name
func (r *TerraformResourceNameContainsTypeRule) Name() string {
	return "terraform_resource_name_contains_type"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformResourceNameContainsTypeRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *TerraformResourceNameContainsTypeRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *TerraformResourceNameContainsTypeRule) Link() string {
	return ""
}

// Check checks whether resource names contain parts of the resource type
func (r *TerraformResourceNameContainsTypeRule) Check(runner tflint.Runner) error {
	content, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "resource",
				LabelNames: []string{"type", "name"},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range content.Blocks {
		resourceType := resource.Labels[0]
		resourceName := resource.Labels[1]
		
		if r.nameContainsType(resourceType, resourceName) {
			err := runner.EmitIssue(
				r,
				fmt.Sprintf("Resource name '%s' contains parts of the resource type '%s'", resourceName, resourceType),
				resource.DefRange,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// nameContainsType checks if the resource name contains any part of the resource type
func (r *TerraformResourceNameContainsTypeRule) nameContainsType(resourceType, resourceName string) bool {
	// Convert both to lowercase for case-insensitive comparison
	lowerType := strings.ToLower(resourceType)
	lowerName := strings.ToLower(resourceName)
	
	// Split resource type by underscores to get individual components
	typeParts := strings.Split(lowerType, "_")
	
	// Check if the full resource type appears in the name
	if strings.Contains(lowerName, lowerType) {
		return true
	}
	
	// Check if any meaningful part of the resource type appears in the name
	// Skip very short parts (like "s3", "ec2") and common prefixes
	for _, part := range typeParts {
		if len(part) > 2 && part != "aws" && part != "gcp" && part != "azurerm" {
			if strings.Contains(lowerName, part) {
				return true
			}
		}
	}
	
	return false
}