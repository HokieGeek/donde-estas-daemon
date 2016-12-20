#!/bin/bash

here=$(cd ${0%/*}; pwd)

host=$(hostname --fqdn)
port=8585

ids=(
    "42"
    "99"
)

dbname="couchdb_$$"
servername="donded_$$"

## Create the database
docker run --detach --volume /usr/local/var/lib/couchdb:/usr/local/var/lib/couchdb --name ${dbname} couchdb:1.6

## Create the server
docker build --tag hokiegeek/donde-estas-daemon ${here}/..
docker run --detach --name ${servername} --publish ${port}:8080 --link ${dbname}:db hokiegeek/donde-estas-daemon

## Add some "real" data
lat="39.18"$(shuf -i 0-9999 -n 1)
lon="-77.27"$(shuf -i 0-9999 -n 1)

${here}/person.sh update ${host}:${port} ${ids[0]} lat:${lat} lon:${lon} white:${ids[1]}
${here}/person.sh update ${host}:${port} ${ids[1]} lat:${lat} lon:${lon} follow:${ids[0]} white:${ids[0]}

## Query that data
${here}/person.sh get ${host}:${port} ${ids[0]}
${here}/person.sh get ${host}:${port} ${ids[1]}

## TODO: perform some sort of data validation :)

## Cleanup
docker stop ${servername} && docker rm ${servername}
docker stop ${dbname} && docker rm ${dbname}
