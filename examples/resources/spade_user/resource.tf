resource "spade_user" "my_user" {
  first_name = "John"
  last_name  = "Doe"
  email      = "john.doe@example.com"

  active = true
  groups = [spade_group.my_group.id]
}
