#!/bin/sh

curl -v `hostname --fqdn`:8585/update --data \
'{
    "id" : "42",
    "name" : "foobar",
    "position": {
        "tov": "2009-11-10T23:00:00Z",
        "latitude": 0.1,
        "longitude": 0.2,
        "elevation": 0.3
    },
    "visible" : true,
    "whitelist" : ["99"],
    "following" : []
}'
curl -v `hostname --fqdn`:8585/update --data \
'{
    "id" : "99",
    "name" : "raboof",
    "position": {
        "tov": "2009-11-10T22:00:00Z",
        "latitude": 0.4,
        "longitude": 0.5,
        "elevation": 0.6
    },
    "visible" : true,
    "whitelist" : ["42"],
    "following" : ["42"]
}'
curl -v `hostname --fqdn`:8585/person --data '{"ids":["42","99"]}'
