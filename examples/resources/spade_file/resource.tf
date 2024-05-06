resource "spade_file" "my_file" {
  code = "File upload"
  format = spade_file_format.my_format.id
  processor = spade_file_processor.my_processor.id
  tags = ["File 1"]
  system_params = jsonencode({
    foo = "bar"
    baz = "qux"
  })
}
