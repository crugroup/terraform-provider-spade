resource "spade_process" "my_process" {
  code = "Run process"
  executor = spade_executor.my_executor.id
  tags = ["File 1"]
  system_params = jsonencode({
    foo = "bar"
    baz = "qux"
  })
}
