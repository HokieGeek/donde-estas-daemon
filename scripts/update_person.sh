#!/bin/sh

host=$1
id=$2

if [ $# -gt 2 ]; then
    lat=$3
    lon=$4
else
    lat="39.18"$(shuf -i 0-9999 -n 1)
    lon="-77.27"$(shuf -i 0-9999 -n 1)
fi

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
    "whitelist" : ["andres","olivia"],
    "following" : ["andres","olivia"]
}'

curl -v ${host}/person --data '{"ids":["'${id}'"]}'
