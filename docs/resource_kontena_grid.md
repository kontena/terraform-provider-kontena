# Resource `kontena_grid`

## Example

```
resource "kontena_grid" "test" {
  name = "test"
  initial_size = 1
  trusted_subnets = [ "192.168.66.0/24" ]
}
```

## Arguments

See [Kontena Platform Configuration Options](https://www.kontena.io/docs/using-kontena/platform.html#configuration-options).

Argument            | Type          |                     | Description
--------------------|---------------|---------------------|-------------
`name`              | String        | Required            | Grid name
`initial_size`      | Int           | Required            | Grid [initial size](https://kontena.io/docs/using-kontena/#initial-nodes)
`token`             | String        | Optional            | Override static token for grid
`subnet`            | String        | Optional            | Override [overlay network subnet](https://kontena.io/docs/advanced/networking.html#subnet-and-supernet)
`supernet`          | String        | Optional            | Override [overlay network supernet](https://kontena.io/docs/advanced/networking.html#subnet-and-supernet)
`default_affinity`  | List<String>  | Optional, Updatable | Default [affinity](https://kontena.io/docs/using-kontena/affinities.html) for services
`trusted_subnets`   | List<String>  | Optional, Updatable | Grid [Trusted Subnets](https://www.kontena.io/docs/advanced/grids.html#manage-kontena-platform-grid-trusted-subnets)

## Attributes

The following arguments are exported as computed attributes:

Attribute   | Type    | Description
------------|---------|--------------------------
`token`     | String  | Generated token for grid
`subnet`    | String  | Default subnet for grid
`supernet`  | String  | Default supernet for grid

### Example

#### `terraform show`
```
kontena_grid.grid:
  id = test
  default_affinity.# = 0
  initial_size = 1
  name = test
  subnet = 10.81.0.0/16
  supernet = 10.80.0.0/12
  token = Gc...
  trusted_subnets.# = 0
```
