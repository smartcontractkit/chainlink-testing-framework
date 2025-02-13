#!/bin/bash
set -x

echo "$CHANGED_FILES"

CHANGED_DIRS=$(
  for FILE in ${CHANGED_FILES}
  do
    # Finds directories with "Dockerfile" files in them
    # in directories that have changed
    DIR=$(dirname "$FILE")
    if [ "$DIR" != "." ]
    then
      find "$DIR" -type d -exec test -e '{}'/Dockerfile \;  -print | sed -e 's/^images\///'
    fi
  done | uniq | jq -c --slurp --raw-input 'split("\n")[:-1] | unique | { "image": . }'
)

echo "dirs=${CHANGED_DIRS}" >> "$GITHUB_OUTPUT"
