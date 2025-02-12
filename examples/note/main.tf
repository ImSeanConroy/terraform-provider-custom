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

resource "custom_note" "example" {
  text = "New Note Example"
}
 
resource "custom_note" "import" {
  text = "Learn Robotics and SLAM"
}

output "example" {
  value = custom_note.example
}