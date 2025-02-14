# Manage example task.
resource "custom_task" "example" {
  title       = "Build custom terraform provider"
  description = "Develop and test a custom terraform provider"
  complete    = false
}
