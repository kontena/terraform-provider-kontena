## Requirements

* [Terraform](https://www.terraform.io/downloads.html) 0.10+
* [Kontena](https://github.com/kontena/kontena) 1.4+

## Install

Install latest release version from GitHub to the terraform [third-party plugins](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins) directory:

    mkdir -p ~/.terraform.d/plugins && curl -L https://gh-releases.kontena.io/kontena/terraform-provider-kontena/gz/latest | tar -C ~/.terraform.d/plugins -xzv

## Development

### Requirements

* [Go](https://golang.org/doc/install) 1.9+

### Build

    go get github.com/kontena/terraform-provider-kontena

### Setup

See [Installing a Plugin](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) from the Terraform docs:

    mkdir -p ~/.terraform.d/plugins && ln -s $GOPATH/bin/terraform-provider-kontena ~/.terraform.d/plugins/

## Documentation

### Provider
* [Provider `kontena`](docs/provider.md)

### Resources
* [Resource `kontena_token`](docs/resource_kontena_token.md)
* [Resource `kontena_grid`](docs/resource_kontena_grid.md)
* [Resource `kontena_node`](docs/resource_kontena_node.md)
* [Resource `kontena_external_registry`](docs/resource_kontena_external_registry.md)

### Example

```
provider "kontena" {
  url = "http://192.168.66.1"
  token = "cc1f...fb6417"
}

resource "kontena_grid" "test" {
  name = "test"
  initial_size = 1
  trusted_subnets = [ "192.168.66.0/24" ]
}

resource "kontena_node" "node" {
  count = "${var.digitalocean_node_count}"

  grid = "${kontena_grid.grid.name}"
  name = "terraform-test-node${count.index + 1}"

  labels = [
    "provider=digitalocean",
    "region=${var.digitalocean-region}",
    "az=${var.digitalocean-region}",
  ]
}

output "grid_token" {
  value = "${kontena_grid.grid.token}"
}
output "node_token" {
  value = "${kontena_node.node.token}"
}
```
