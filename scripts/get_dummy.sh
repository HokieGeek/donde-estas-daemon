#!/bin/bash

host=$1; shift

ids=()
while (( $# > 0 )); do
    (( ${#ids[@]} > 0 )) && delim="," || delim=""
    ids+=(${delim}'"'${1}'"')
    # ids+=('"'${1}'",')
    shift
done

# echo ${ids[@]}
# exit 42

set -x
curl -v ${host}/person --data '{"ids":['"${ids[*]}"']}'
# curl -v ${host}/person --data '{"ids":["andres","keri","olivia"]}'
