# Resource `kontena_external_registry`

## Example

```
resource "kontena_external_registry" "images" {
  grid = "${kontena_grid.grid.name}"
  url         = "https://kontena-images.example.com"
  username    = "test"
  email       = "test@example.com"
  password    = "password"
}
```

## Arguments

Argument        | Type    |                 | Description
----------------|---------|-----------------|-------------
`url`           | String  | Required        | `http://` or `https://` URL
`username`      | String  | Required        | Login username
`password`      | String  | Optional        | Login password
`email`         | String  | Optional (1.5+) |

## Attributes

The following computed attributes are exported:

Attribute         | Type        | Description
------------------|-------------|--------------------------
`id`              | String      | `:grid/:name`
`name`            | String      | URL hostname

### Example

#### `terraform show`
```
kontena_external_registry.test:
  id = demo-3/kontena-images.example.com
  email = test@example.com
  grid = demo-3
  name = kontena-images.example.com
  password = password
  url = https://kontena-images.example.com
  username = test
```

## Importing

Existing `kontena_external_registry` resources can be imported using `terraform import kontena_external_registry.NAME <grid>/<name>`: `terraform import kontena_external_registry.test test/images.example.com`

When importing, the `password` argument can be left as an empty string (`password = ""`). Otherwise terraform will forcibly re-create the resource to ensure that the password matches, because it cannot be read back from the API.
