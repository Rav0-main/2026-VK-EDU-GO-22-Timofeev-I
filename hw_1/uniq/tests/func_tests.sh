#!/bin/bash

SCRIPT_DIR="$(realpath $(dirname $0))"
TEST_DIR="$SCRIPT_DIR/data"

cd "$(dirname $SCRIPT_DIR)" || exit 1

echo "POSITIVE TESTS:"

pos_completed=0
pos_count=0
for f in "$TEST_DIR"/*.txt; do
  if [ ! -f "$f" ]; then
    continue
  fi

  fname="$(basename "$f")"
  if printf "%s""$fname" | grep -oE "pos_[0-9]+_in.txt" >/dev/null 2>&1; then
    num=$(printf "%s" "$fname" | grep -oE "[0-9]+")
    fin="$fname"
    fout=pos_"$num"_out.txt
    fargs=pos_"$num"_args.txt

    if [ ! -f "$TEST_DIR"/"$fout" ]; then
      echo Error. For in: \""$fin"\" not found out: \""$fout"\"

    else
      ((pos_count++))

      echo -n "# $pos_count. "
      if "$SCRIPT_DIR/pos_case.sh" "$TEST_DIR/$fin" "$TEST_DIR/$fout" "$TEST_DIR/$fargs"; then
        echo "Test $fin passed."
        ((pos_completed++))
      else
        echo "Test $fin failed."
      fi
    fi
  fi
done

if [ $pos_count -eq 0 ]; then
  echo "Not exists."
fi

echo
echo "============================="
echo "NEGATIVE TESTS:"

neg_completed=0
neg_count=0
for f in "$TEST_DIR"/*.txt; do
  if [ ! -f "$f" ]; then
    continue
  fi

  fname=$(basename "$f")
  if printf "%s" "$fname" | grep -oE "neg_[0-9]+_in.txt" >/dev/null 2>&1; then
    num=$(printf "%s" "$fname" | grep -oE "[0-9]+")
    fin="$fname"
    fargs=neg_"$num"_args.txt

    ((neg_count++))

    echo -n "# $neg_count. "
    if "$SCRIPT_DIR/neg_case.sh" "$TEST_DIR/$fin" "$TEST_DIR/$fargs"; then
      echo "Test $fin passed."
      ((neg_completed++))

    else
      echo "Test $fin failed."
    fi
  fi
done

if [ $neg_count -eq 0 ]; then
  echo "Not exists."
fi

echo
echo Positive: "$pos_completed" / "$pos_count".
echo Negative: "$neg_completed" / "$neg_count".

if [ $((pos_count + neg_count - pos_completed - neg_completed)) -eq 0 ]; then
  echo "PASSED."

else
  echo "FAILED."
fi

exit $((pos_count + neg_count - pos_completed - neg_completed))
