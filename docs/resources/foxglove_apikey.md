---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "foxglove_apikey Resource - terraform-provider-foxglove-cloud"
subcategory: ""
description: |-
   Create and manage foxglove api keys
---

# foxglove_apikey (Resource)

This resource allows you to create and manage [api keys in Foxglove Cloud](https://docs.foxglove.dev/docs/api/#api-keys).

#### Example Usage

```terraform
resource "foxglove_apikey" "foo" {
  label = "Foo api key"
  capabilities = [
    "events.list",
    "events.create",
    "events.update",
    "events.delete",
  ]
}
```

#### Schema

##### Required

- `label` (String) The human-readable label for this key.
- `capabilities` (List of strings) Capabilities of this key

##### Read-Only

- `id` (String) The unique identifier.
- `secret` (String, Sensitive) The secret token.

## Import

Import is not supported at the moment.
