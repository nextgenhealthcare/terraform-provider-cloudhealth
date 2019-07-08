#!/usr/bin/env bash

echo "==> Building..."

if [[ -v VERSION ]]; then
  OUTPUT="pkg/{{.OS}}_{{.Arch}}/{{.Dir}}_${VERSION}"
else
  OUTPUT="pkg/{{.OS}}_{{.Arch}}/{{.Dir}}"
fi
gox -os "freebsd darwin linux windows" -arch "386 amd64" -output "${OUTPUT}"

if [[ -v VERSION ]]; then
  echo "==> Packaging"
  PACKAGE="${PWD##*/}"
  for PLATFORM in $(find ./pkg -mindepth 1 -maxdepth 1 -type d); do
    OSARCH=$(basename ${PLATFORM})
    echo "--> ${OSARCH}"
    pushd $PLATFORM >/dev/null 2>&1
    zip ../${PACKAGE}_${OSARCH}.zip ./*
    popd >/dev/null 2>&1
  done
fi
