# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is a TFLint ruleset plugin that provides custom linting rules for Terraform configurations. It's built using the TFLint Plugin SDK and follows the standard TFLint plugin architecture.

The plugin operates as an independent binary that communicates with TFLint via gRPC. TFLint executes the plugin when enabled and the plugin acts as a gRPC server to provide rule checking functionality.

## Common Commands

### Build and Development
- `make` or `make build` - Build the plugin binary
- `make test` - Run all tests with `go test ./...`
- `make install` - Build and install the plugin to `~/.tflint.d/plugins`
- `go test ./rules` - Run tests for a specific package
- `go test ./rules -run TestName` - Run a specific test

### Testing the Plugin
After building and installing:
```bash
# Create a basic .tflint.hcl config
cat << EOS > .tflint.hcl
plugin "template" {
  enabled = true
}
EOS

# Run TFLint with the plugin
tflint
```

## Architecture

### Plugin Structure
- `main.go` - Plugin entry point that registers all rules with the TFLint plugin system
- `rules/` - Contains all rule implementations and their corresponding tests

### Rule Implementation Pattern
Each rule follows this standard structure:
1. **Rule struct** - Embeds `tflint.DefaultRule` and implements the `tflint.Rule` interface
2. **Constructor function** - `NewRuleName()` returns a new instance
3. **Required methods**:
   - `Name() string` - Returns the rule identifier
   - `Enabled() bool` - Whether enabled by default
   - `Severity() tflint.Severity` - Rule severity (ERROR, WARNING, NOTICE)
   - `Link() string` - Documentation link (can be empty)
   - `Check(runner tflint.Runner) error` - Main rule logic

### Test Pattern
Each rule has a corresponding `*_test.go` file that:
- Uses `helper.TestRunner()` to create test environments
- Defines test cases with Terraform configuration strings
- Uses `helper.AssertIssues()` to verify expected rule violations

### Rule Types Examples
- **aws_instance_example_type** - Shows how to access top-level resource attributes
- **aws_s3_bucket_example_lifecycle_rule** - Demonstrates accessing nested blocks and attributes
- **google_compute_ssl_policy** - Example with custom rule configuration
- **terraform_backend_type** - Shows how to access non-resource configurations

## Adding New Rules

1. Create rule file in `rules/` directory following the naming pattern
2. Implement the required interface methods
3. Create corresponding test file with `_test.go` suffix
4. Register the rule in `main.go` by adding it to the `Rules` slice
5. Use `runner.GetResourceContent()` to access Terraform configuration
6. Use `runner.EvaluateExpr()` to evaluate expressions and emit issues

## Plugin Development Guidelines

### Repository Structure
- Repository name should follow `tflint-ruleset-*` pattern
- Use the `tflint-ruleset` GitHub topic for discoverability
- Based on the official [tflint-ruleset-template](https://github.com/terraform-linters/tflint-ruleset-template)

### Development Workflow
1. Make changes to rules or add new rules
2. Test locally with `make install` followed by running TFLint
3. Use `TFLINT_LOG=debug` environment variable for debug logging
4. Refer to the [TFLint Plugin SDK API reference](https://github.com/terraform-linters/tflint-plugin-sdk) for implementation details

### Release Process
When ready to release:
1. Create version tags like `v1.1.1`
2. Publish binaries on GitHub Releases
3. Include platform-specific zip assets
4. Provide a `checksums.txt` file
5. Optional: Add PGP signing or artifact attestation

## Dependencies

The plugin uses:
- `github.com/terraform-linters/tflint-plugin-sdk` - Core plugin SDK for gRPC communication with TFLint
- `github.com/hashicorp/hcl/v2` - HCL parsing and evaluation