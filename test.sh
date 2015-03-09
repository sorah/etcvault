#!/bin/bash
set -e

PKGS="./keys ./container ./engine"
FORMATS="$PKGS *.go"

for pkg in $PKGS; do
  go test -cover $pkg
done

fmt_result="$(gofmt -l $FORMATS)"
if [ -n "${fmt_result}" ]; then
	echo -e "gofmt checking failed:\n${fmt_result}"
  exit 1
fi

