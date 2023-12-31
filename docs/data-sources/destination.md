---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "opensearch_destination Data Source - terraform-provider-opensearch"
subcategory: ""
description: |-
  opensearch_destination can be used to retrieve the destination object by name.
---

# opensearch_destination (Data Source)

`opensearch_destination` can be used to retrieve the destination object by name.

## Example Usage

```terraform
# Example destination in other terraform plan
# resource "opensearch_destination" "test" {
#   body = <<EOF
# {
#   "name": "my-destination",
#   "type": "slack",
#   "slack": {
#     "url": "http://www.example.com"
#   }
# }
# EOF
# }

data "opensearch_destination" "test" {
  name = "my-destination"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the destrination to retrieve

### Read-Only

- `body` (Map of String) Map of the attributes of the destination
- `id` (String) The ID of this resource.


