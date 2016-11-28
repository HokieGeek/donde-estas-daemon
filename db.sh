#!/bin/bash

name="dondedb"

if [[ $(docker ps --quiet --filter="name=${name}") == "" ]]; then
    echo "Creating an instance of the database"
    # docker run --detach --name ${name} --env COUCHDB_USER=admin --env COUCHDB_PASSWORD=password couchdb
    docker run --detach --name ${name} couchdb
fi

cmd=$1; shift

# set -x
docker run -it --rm --link ${name}:db hokiegeek/alpine-curl -v -X ${cmd:-"GET"} http://db:5984/$@
# docker run -it --rm --link ${name}:db hokiegeek/alpine-curl $@
