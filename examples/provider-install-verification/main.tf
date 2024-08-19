terraform {
  required_providers {
    custom = {
      source = "hashicorp.com/edu/custom"
    }
  }
}

provider "custom" {}

data "custom_coffees" "example" {}
