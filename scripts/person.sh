#!/bin/bash

cmd=$1
host=$2
shift 2

case ${cmd} in
update)
    id=$1; shift
    lat="39.18"$(shuf -i 0-9999 -n 1)
    lon="-77.27"$(shuf -i 0-9999 -n 1)
    following=""
    whitelist=""
    while (( $# > 0 )); do
        key=${1%%:*}
        val=${1##*:}
        case ${key} in
        lat) lat=${val} ;;
        lon) lon=${val} ;;
        follow) following="\"${val//,/\",\"}\"" ;;
        white) whitelist="\"${val//,/\",\"}\"" ;;
        esac
        shift
    done

    curl -v ${host}/update --data \
    '{
        "id" : "'${id}'",
        "name" : "'${id}'",
        "position": {
            "tov": "'$(date +"%FT%TZ")'",
            "latitude": '${lat}',
            "longitude": '${lon}',
            "elevation": 0.0
        },
        "visible" : true,
        "whitelist" : ['${whitelist}'],
        "following" : ['${following}']
    }'
    ;;
get)
    ids=()
    while (( $# > 0 )); do
        (( ${#ids[@]} > 0 )) && delim="," || delim=""
        ids+=(${delim}'"'${1}'"')
        shift
    done

    curl ${host}/person --data '{"ids":['"${ids[*]}"']}'
    ;;
esac
