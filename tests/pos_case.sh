#!/bin/bash

SRC_DIR="$(dirname $(dirname $0))"
COMPARATOR="$(dirname $0)/comparator.sh"

APP="$SRC_DIR/uniq"
OUT_FNAME=".uniq.test.tmp.out"

if [ $# != 3 ]; then
  echo Error. Need write: "$0" \"file_in\" \"file_expecting_out\" \"file_args\".
  exit 1
fi

fin="$1"
fout="$2"
fargs="$3"

if [ ! -f "$fin" ]; then
  echo Error. \""$fin"\" not file.
  exit 1

elif [ ! -f "$fout" ]; then
  echo Error. \""$fout"\" not file.
  exit 1

elif [ "$OUT_FNAME" == "$fin" ] || [ "$OUT_FNAME" == "$fout" ]; then
  echo Error. Don\'t use \""$OUT_FNAME"\" for in-file or out-file.
  exit 1

elif [ ! -f "$APP" ]; then
  echo Error. Not found compiled app \""$APP"\".
  exit 1
fi

args=""
if [ -f "$fargs" ]; then
  args="$(cat $fargs 2>/dev/null)"
fi

echo "Result of: $(basename $APP) $args < $(basename $fin) => $(basename $fout):"
if eval "$APP $(cat $fargs 2>/dev/null)" <"$fin" >"$SRC_DIR/$OUT_FNAME" 2>&1; then
  if "$COMPARATOR" "$SRC_DIR/$OUT_FNAME" "$fout"; then
    echo "OK!"
    exit 0
  else
    echo "- WA. Expected:"
    cat -n "$fout"
    echo
    echo "- But given:"
    cat -n "$SRC_DIR/$OUT_FNAME"
    exit 1
  fi
else
  echo RE. Ended with error code: $?
  exit 1
fi
