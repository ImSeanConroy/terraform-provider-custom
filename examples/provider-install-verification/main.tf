terraform {
  required_providers {
    custom = {
      source = "hashicorp.com/edu/custom"
    }
  }
}

data "custom_tasks" "example" {}
