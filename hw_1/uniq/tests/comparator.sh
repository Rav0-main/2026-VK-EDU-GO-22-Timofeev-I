#!/usr/bin/bash

if [ $# != 2 ]; then
  echo "Error. Need write $(basename $0) 'first_file' 'second_file'"
fi

diff $1 $2 >/dev/null 2>&1
