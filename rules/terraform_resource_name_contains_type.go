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
		
		if r.hasSharedSuffix(resourceType, resourceName) {
			err := runner.EmitIssue(
				r,
				fmt.Sprintf("Resource name '%s' should not have a shared suffix with the resource type '%s'.", resourceName, resourceType),
				resource.DefRange,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// hasSharedSuffix checks if the resource name has a shared suffix with the resource type
func (r *TerraformResourceNameContainsTypeRule) hasSharedSuffix(resourceType, resourceName string) bool {
	// Convert both to lowercase for case-insensitive comparison
	lowerType := strings.ToLower(resourceType)
	lowerName := strings.ToLower(resourceName)
	
	// Split both by underscores to get word segments
	typeParts := strings.Split(lowerType, "_")
	nameParts := strings.Split(lowerName, "_")
	
	// Check all possible suffixes of the resource type
	for i := 0; i < len(typeParts); i++ {
		// Get suffix starting from position i
		typeSuffix := typeParts[i:]
		
		// Check if name ends with this suffix
		if len(nameParts) >= len(typeSuffix) {
			nameStart := len(nameParts) - len(typeSuffix)
			nameSuffix := nameParts[nameStart:]
			
			// Compare the suffixes
			match := true
			for j := 0; j < len(typeSuffix); j++ {
				if typeSuffix[j] != nameSuffix[j] {
					match = false
					break
				}
			}
			
			if match {
				return true
			}
		}
	}
	
	return false
}