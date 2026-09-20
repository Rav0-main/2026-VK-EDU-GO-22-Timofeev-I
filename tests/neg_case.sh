#!/bin/bash

SRC_DIR="$(dirname $(dirname $0))"
APP="$SRC_DIR/uniq"
OUT_FNAME=".uniq.test.tmp.out"

if [ $# != 2 ]; then
  echo Error. Need write: "$(basename $0)" \"file_in\" \"file_args\"
  exit 1
fi

fin="$1"
fargs="$2"

if [ ! -f "$fin" ]; then
  echo Error. \""$fin"\" not file.
  exit 1

elif [ ! -f "$APP" ]; then
  echo Error. Not found compiled app \""$APP"\".
  exit 1

fi

args=""
if [ -f "$fargs" ]; then
  args="$(cat $fargs 2>/dev/null)"
fi

echo "Result of: $(basename $APP) $args < $(basename $fin) => ERROR:"
if ! eval "$APP $(cat $fargs) < $fin > $SRC_DIR/$OUT_FNAME" 2>&1; then
  echo "OK!"
  exit 0
else
  echo Error. Ended with success "(return code = 0)".
  exit 1
fi
