---
page_title: "Foxglove Cloud Provider"
subcategory: ""
description: |-
    The Foxglove Cloud provider provides resources to manage resources in Foxglove Cloud.
---

# Foxglove Cloud Provider

The Foxglove Cloud provider uses the [Foxglove Cloud api](https://docs.foxglove.dev/docs/api/) to manage devices and api keys.

The changelog for this provider can be found here: <https://github.com/siinm/terraform-provider-foxglove-cloud/releases>.

This provider is independently developed and is not affiliated with the Foxglove company.

## Example Usage

### Creating a Foxglove Cloud provider

```terraform
provider "foxglove" {
  api_key    = "abc"
}

resource "foxglove_device" "foo" {
  name = "foo"
}
```

## Schema

### Optional

- `api_key` (String) Foxglove API Key. Can also be set via environment variable FOXGLOVE_API_KEY

## Functions

Currently, the Foxglove Cloud provider does not support any functions.

## Data Sources

Currently, the Foxglove Cloud provider does not support any data sources.
