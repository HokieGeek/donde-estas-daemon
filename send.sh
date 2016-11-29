#!/bin/sh

curl -v `hostname --fqdn`:8585/update --data \
'{
    "id" : 42,
    "name" : "foobar",
	"position": {
		"tov": "2009-11-10T23:00:00Z",
		"latitude": 0.0,
		"longitude": 0.0,
		"elevation": 0.0
	},
	"visible" : true,
	"whitelist" : [],
	"following" : []
}'
curl -v `hostname --fqdn`:8585/update --data \
'{
    "id" : 99,
    "name" : "raboof",
	"position": {
		"tov": "2009-11-10T23:00:00Z",
		"latitude": 0.0,
		"longitude": 0.0,
		"elevation": 0.0
	},
	"visible" : true,
	"whitelist" : [],
	"following" : []
}'
curl -v `hostname --fqdn`:8585/person --data '{"ids":[42,99]}'
