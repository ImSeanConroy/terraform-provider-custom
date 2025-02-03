terraform {
  required_providers {
    custom = {
      source = "hashicorp.com/edu/custom"
    }
  }  
}

provider "custom" {
  host = "http://localhost:3000"
  token = "21097dc18cfb54ebb691098b6dcf556418f9a6d2640db2217a529c"
}

data "custom_notes" "example" {}
 