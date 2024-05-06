resource "spade_file_processor" "my_processor" {
  name = "My file processor"
  callable = "spadeapp.examples.processor.ExampleFileProcessor"
}
