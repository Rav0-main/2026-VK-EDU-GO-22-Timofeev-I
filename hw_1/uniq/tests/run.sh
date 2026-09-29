#!/bin/bash

SCRIPT_DIR="$(dirname $(realpath $0))"

cd "$(dirname $SCRIPT_DIR)"

# build application
if ! go build -o uniq.out .; then
  echo "Error. Can not build 'uniq.out'."
  exit 1
fi

echo
cd "$SCRIPT_DIR"

./func_tests.sh
