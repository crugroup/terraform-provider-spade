resource "spade_executor" "my_executor" {
  name                      = "My executor"
  callable                  = "spadeapp.examples.executor.ExampleExecutor"
  history_provider_callable = "spadeapp.examples.executor.ExampleHistoryProvider"
}
