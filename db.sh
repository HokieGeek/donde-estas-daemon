#!/bin/bash

name="dondedb"

if [[ $(docker ps --quiet --filter="name=${name}") == "" ]]; then
    echo "Creating an instance of the database"
    # docker run --detach --volume /usr/local/var/lib/couchdb:/usr/local/var/lib/couchdb --name dondedb --env COUCHDB_USER=donde --env COUCHDB_PASSWORD=dondepwd couchdb
    docker run --detach --volume /usr/local/var/lib/couchdb:/usr/local/var/lib/couchdb --name ${name} couchdb
fi

# cmd=$1; shift

# set -x
# docker run -it --rm --link ${name}:db hokiegeek/alpine-curl -v -X ${cmd:-"GET"} http://db:5984/$@
# docker run -it --rm --link ${name}:db hokiegeek/alpine-curl $@
