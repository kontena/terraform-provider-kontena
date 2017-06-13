## Example

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
```
