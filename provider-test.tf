terraform {
  required_providers {
    spade = {
      source = "crugroup/spade"
    }
  }
}

provider "spade" {
  url = "https://spade-backend-dev.crugroup.com"
  email = "spade@crugroup.com"
  password = "spadespadespade"
}

resource "spade_executor" "dummy_executor" {
  name = "Dummy Executor2"
  callable = "spadeapp.examples.executor.ExampleExecutor"
}

resource "spade_file_processor" "dummy_processor" {
  name = "Dummy Processor2"
  callable = "spadeapp.examples.processor.ExampleFileProcessor"
}

resource "spade_file_format" "dummy_format" {
  format = "json2"
}

resource "spade_process" "test_process" {
  code = "test2"
  executor = spade_executor.dummy_executor.id
  description = "hello world"
  tags = ["Steel"]
  system_params = jsonencode({
    a = "b"
  })
}

resource "spade_file" "test_file" {
  code = "test2"
  format = spade_file_format.dummy_format.id
  processor = spade_file_processor.dummy_processor.id
  description = "hello world"
  tags = ["Steel"]
  system_params = jsonencode({
    a = "b"
  })
}
