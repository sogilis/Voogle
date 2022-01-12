#!/usr/bin/env bash

DOCKERFILES=$(git ls-files | grep '\Dockerfile$')

# COLOR
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

STATUS=0

for file in $DOCKERFILES
do
    echo -e "$GREEN Checking -> $file $NC"
    docker run --rm -i hadolint/hadolint < "$file"
    # shellcheck disable=SC2181
    if [ $? -ne 0 ]
    then
       echo -e "$RED FAILED check on $file $NC"
       STATUS=1
    fi
done

exit $STATUS
