package main

import (
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/cinder-technologies/tflint-ruleset-cinder/rules"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "cinder",
			Version: "0.1.0",
			Rules: []tflint.Rule{
				rules.NewTerraformResourceNameContainsTypeRule(),
			},
		},
	})
}
