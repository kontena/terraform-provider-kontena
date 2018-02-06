#!/bin/bash

set -uex

ARCH=linux_amd64
PROVIDER=kontena

mkdir -p build/dist

go build -o build/dist/terraform-provider-${PROVIDER}_v${VERSION} .

( cd build/dist

  zip terraform-provider-${PROVIDER}_${VERSION}_${ARCH}.zip terraform-provider-${PROVIDER}_v${VERSION}
  tar -czvf terraform-provider-${PROVIDER}_${VERSION}_${ARCH}.tar.gz terraform-provider-${PROVIDER}_v${VERSION}

  sha256sum *.tar.gz *.zip > SHA256SUM
)
