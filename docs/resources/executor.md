---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "spade_executor Resource - spade"
subcategory: ""
description: |-
  Represents an executor within Spade
---

# spade_executor (Resource)

Represents an executor within Spade

## Example Usage

```terraform
resource "spade_executor" "my_executor" {
  name                      = "My executor"
  callable                  = "spadeapp.examples.executor.ExampleExecutor"
  history_provider_callable = "spadeapp.examples.executor.ExampleHistoryProvider"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `callable` (String) Python import path to the Executor class
- `name` (String) Name of the executor

### Optional

- `description` (String) Description of the executor
- `history_provider_callable` (String) Python import path to the HistoryProvider class

### Read-Only

- `id` (Number) Identifier of the executor
