# Provider `kontena`

## Example

```
provider "kontena" {
  url = "http://192.168.66.1:9292"
  token = "ee7a431f238c2681154d5f066236ea08b14abbefc8306bb86736d635017e6e3a"
}
```

## Arguments

Argument        | Type   |          | Description
----------------+--------+----------+------
`url`           | String | Required | `http://` or `https://` URL
`token`         | String | Optional | Access token
`ssl_cert_pem`  | String | Optional | Validate `https` certificate using PEM-encoded CA cert
`ssl_cert_cn`   | String | Optional | Override `https` SNI hostname, cerificate validation Common Name

### Access `token`

The access token can be generated using `kontena master token create` for an existing Kontena Master.

### Bootstrapping the access token using the  `INITIAL_ADMIN_CODE`

If you are provisioning the master itself via terraform with a pre-configured `INITIAL_ADMIN_CODE`, you can use the `kontena_token` resource to exchange the inital admin code for the admin token. You must use a second, aliased instance of the provider to do this:

```
variable "master_url" {

}

variable "initial_admin_code" {

}


provider "kontena" {
  alias = "master-bootstrap"
  url = "${var.master_url}"
}

resource "kontena_token" "admin" {
  provider = "kontena.master-bootstrap"

  code = "${var.initial_admin_code}"
}

provider "kontena" {
  url = "${module.digitalocean_master.http_url}"
  token = "${kontena_token.admin.token}"
}

output "KONTENA_URI" {
  value = "${var.master_url}"
}

output "KONTENA_TOKEN" {
  value = "${kontena_token.admin.token}"
}
```
