resource "spade_variable_set" "my_variable_set" {
  name        = "my_variable_set"
  description = "..."
  
  variables = [
    spade_variable.my_variable.id,
    spade_secret_variable.my_secret.id
  ]
}
