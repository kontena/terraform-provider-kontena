#!/bin/bash

set -uex

PROVIDER=kontena

mkdir -p build/dist

function build_arch() {
  local os=$1
  local arch=$2

  mkdir -p build/dist-${os}_${arch}

  GOOS=$os GOARCH=$arch go build -o build/dist-${os}_${arch}/terraform-provider-${PROVIDER}_v${VERSION} .

  ( cd build/dist-${os}_${arch}

    zip ../dist/terraform-provider-${PROVIDER}_${VERSION}_${os}_${arch}.zip terraform-provider-${PROVIDER}_v${VERSION}
    tar -czvf ../dist/terraform-provider-${PROVIDER}_${VERSION}_${os}_${arch}.tar.gz terraform-provider-${PROVIDER}_v${VERSION}
  )
}

build_arch linux amd64
build_arch darwin amd64
build_arch windows amd64

( cd build/dist

  sha256sum *.tar.gz *.zip > SHA256SUM
)
