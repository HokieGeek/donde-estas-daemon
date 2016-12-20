#!/bin/bash

host=$1; shift

ids=()
while (( $# > 0 )); do
    (( ${#ids[@]} > 0 )) && delim="," || delim=""
    ids+=(${delim}'"'${1}'"')
    shift
done

curl ${host}/person --data '{"ids":['"${ids[*]}"']}'
