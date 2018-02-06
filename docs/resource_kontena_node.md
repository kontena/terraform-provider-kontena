# Resource `kontena_node`

## Example

```
resource "kontena_node" "node" {
  grid = "${kontena_grid.grid.name}"
  name = "terraform-test-node${var.index}"

  labels = [
    "provider=${var.provider}",
    "region=${var.region}",
    "az=${var.az}",
  ]
}
```

## Arguments

Argument        | Type          |                     | Description
----------------+---------------+---------------------+-------------
`grid`          | String        | Required            | Grid name
`name`          | String        | Required            | Node name (unique)
`token`         | String        | Optional, Updatable | Node token (if not using Grid token)
`labels`        | List<String>  | Optional, Updatable | Node labels


## Attributes

The following computed attributes are exported:

Attribute         | Type        | Description
------------------+-------------+--------------------------
`node_id`         | String      | Docker ID
`node_number`     | Integer     | Unique node number within grid
`initial_node`    | Boolean     | Node is a [Grid Initial Node](https://kontena.io/docs/using-kontena/#initial-nodes)
`public_ip`       | String      |
`private_ip`      | String      |
`overlay_ip`      | String      |
`agent_version`   | String      |
`docker_version`  | String      |
