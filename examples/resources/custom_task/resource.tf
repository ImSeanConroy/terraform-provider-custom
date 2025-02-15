# Manage example task.
resource "custom_task" "example" {
  title       = "Build terraform resource"
  description = "Develop and test a custom terraform resource"
  complete    = false
}
