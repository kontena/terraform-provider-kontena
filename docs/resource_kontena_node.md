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

Argument        | Type     |                     | Description
----------------|----------|---------------------|-------------
`grid`          | String   | Required            | Grid name
`name`          | String   | Required            | Node name (unique)
`token`         | String   | Optional, Updatable | Node token (if not using Grid token)
`labels`        | [String] | Optional, Updatable | Node labels


## Attributes

The following computed attributes are exported:

Attribute         | Type        | Description
------------------|-------------|--------------------------
`node_id`         | String      | Docker ID
`node_number`     | Integer     | Unique node number within grid
`initial_node`    | Boolean     | Node is a [Grid Initial Node](https://kontena.io/docs/using-kontena/#initial-nodes)
`public_ip`       | String      |
`private_ip`      | String      |
`overlay_ip`      | String      |
`agent_version`   | String      |
`docker_version`  | String      |

### Example

#### `terraform show`
```
module.digitalocean-node2.kontena_node.node:
  id = test/terraform-test-node2
  agent_version = 1.4.0.pre5
  docker_version = 1.12.6
  grid = test
  initial_node = true
  labels.# = 3
  labels.0 = provider=digitalocean
  labels.1 = region=fra1
  labels.2 = az=fra1
  name = terraform-test-node2
  node_id = PBOE:X2WG:ULOL:O2YJ:6C3L:P32Z:P2Y3:EGAR:27AJ:NMXG:Z2SU:TWFZ
  node_number = 1
  overlay_ip = 10.81.0.1
  private_ip = 46.101.119.165
  public_ip = 46.101.119.165
  token = cFdn...lQ==
module.digitalocean-node1.kontena_node.node:
  id = test/terraform-test-node1
  agent_version = 1.4.0.pre5
  docker_version = 1.12.6
  grid = test
  initial_node = false
  labels.# = 3
  labels.0 = provider=digitalocean
  labels.1 = region=fra1
  labels.2 = az=fra1
  name = terraform-test-node1
  node_id = AHBV:DHKM:5J3Z:WST3:DO2Q:Q2AG:T6NM:BBYB:QGN5:G4Y6:7GXZ:4QIT
  node_number = 2
  overlay_ip = 10.81.0.2
  private_ip = 207.154.200.134
  public_ip = 207.154.200.134
  token = cNr...343g==
```

## Importing

Existing `kontena_node` resources can be imported using `terraform import kontena_node.NAME <grid>/<name>`: `terraform import kontena_node.node-1 test/node-1`
