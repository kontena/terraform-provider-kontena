## Example

#### `example.tf`
```
provider "kontena" {
  url = "http://${digitalocean_droplet.master.ipv4_address}:80"
  access_token = "${var.kontena-access_token}"
}

resource "kontena_grid" "test" {
  name = "test"
  initial_size = 1
  trusted_subnets = [ "192.168.66.0/24" ]
}

resource "kontena_node" "node" {
  count = "${var.digitalocean_node_count}"
  depends_on = [
    "kontena_grid.grid",
    "digitalocean_droplet.node",
  ]

  grid = "${kontena_grid.grid.name}"
  name = "terraform-test-node${count.index + 1}"

  labels = [
    "provider=digitalocean",
    "region=${digitalocean_droplet.node.*.region[count.index]}",
    "az=${digitalocean_droplet.node.*.region[count.index]}",
  ]
}
```

#### `terraform show`
```
kontena_grid.test:
  id = test
  default_affinity.# = 0
  initial_size = 2
  name = test
  subnet = 10.81.0.0/16
  supernet = 10.80.0.0/12
  token = 8TyJCGXvHGLVfOIN/44RyeqEQs1MTMiDVTO/HJjQ+bJHrM7ZGOMZdzpbJyKS8P6ObU4xxf/M6hM8vRCqO3OUOQ==
  trusted_subnets.# = 1
  trusted_subnets.0 = 192.168.66.0/24
kontena_node.node:
  id = test/terraform-test-node1
  agent_version = 1.3.0
  docker_version = 1.12.6
  grid = test
  initial_node = true
  labels.# = 3
  labels.0 = provider=digitalocean
  labels.1 = region=fra1
  labels.2 = az=fra1
  name = terraform-test-node1
  node_number = 1
  overlay_ip = 10.81.0.1
  private_ip = 207.154.198.57
  public_ip = 207.154.198.57
```
