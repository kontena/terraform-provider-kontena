# Resource `kontena_token`

## Example

```
resource "kontena_token" "admin" {
  provider = "kontena.master-bootstrap"

  code = "${var.initial_admin_code}"
}
```

## Arguments

Argument        | Type    |           | Description
----------------|---------|-----------|-------------
`code`          | String  | Optional  | OAuth2 code (also `INITIAL_ADMIN_CODE`)

## Attributes

The following computed attributes are exported:

Attribute   | Type      | Description
------------|-----------|--------------------------
`token`     | String    | OAuth2 token for `provider "kontena"`
`user`      | String    | Kontena Master user associated with token
`email`     | String    | Kontena Cloud login email for user
`roles`     | [String]  | Kontena Master user roles for token (e.g. `master_admin`)

### Example

#### `terraform show`
```
kontena_token.admin:
  id = 971...
  code = bsF...
  email = admin
  roles.# = 1
  roles.0 = master_admin
  token = 971...
  user = admin
```
