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
  trusted_subnets.# = 0
```
