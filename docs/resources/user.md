---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "onepasswordorg_user Resource - terraform-provider-onepasswordorg"
subcategory: ""
description: |-
  Provides a User resource.
  When a 1password user resources is created, it will be invited  by email.
---

# onepasswordorg_user (Resource)

Provides a User resource.

When a 1password user resources is created, it will be invited  by email.

## Example Usage

```terraform
resource "onepasswordorg_user" "user0" {
  name  = "User zero"
  email = "user0@slok.dev"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `email` (String) The description of the user.
- `name` (String) The name of the user.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
# Go to the website and get the UUID from the URL or use the `op` cli:
op user get user0@slok.dev

# Import.
terraform import onepasswordorg_user.user0 ${ONEPASSWORD_UUID}
```
