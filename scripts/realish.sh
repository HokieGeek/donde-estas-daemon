#!/bin/sh

host=$1

curl -v ${host}/update --data \
'{
    "id" : "keri",
    "name" : "keri",
    "position": {
        "tov": "2009-11-10T22:00:00Z",
        "latitude": 39.189658,
        "longitude": -77.279528,
        "elevation": 0.0
    },
    "visible" : true,
    "whitelist" : ["andres","olivia"],
    "following" : ["andres","olivia"]
}'
curl -v ${host}/update --data \
'{
    "id" : "olivia",
    "name" : "olivia",
    "position": {
        "tov": "2009-11-10T22:00:00Z",
        "latitude": 39.1888622,
        "longitude": -77.287454,
        "elevation": 0.0
    },
    "visible" : true,
    "whitelist" : ["andres","keri"],
    "following" : ["andres","keri"]
}'
curl -v ${host}/update --data \
'{
    "id" : "andres",
    "name" : "foobar",
    "position": {
        "tov": "2009-11-10T23:00:00.210Z",
        "latitude": 0,
        "longitude": 0,
        "elevation": 0
    },
    "visible" : true,
    "whitelist" : ["keri","olivia"],
    "following" : ["keri","olivia"]
}'

curl -v ${host}/person --data '{"ids":["andres","keri","olivia"]}'
