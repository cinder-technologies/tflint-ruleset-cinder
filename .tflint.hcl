plugin "cinder" {
  enabled = true
  source  = "github.com/cinder-technologies/tflint-ruleset-cinder"
  version = "0.1.0"
}

rule "terraform_resource_name_contains_type" {
  enabled = true
}