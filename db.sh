#!/bin/bash

name="dondedb"

if [[ $1 == "--kill" ]]; then
    docker stop ${name}
    docker kill ${name}
    docker rm ${name}
else
    if [[ $(docker ps --quiet --filter="name=${name}") == "" ]]; then
        echo "Creating an instance of the database"
        # docker run --detach --volume /usr/local/var/lib/couchdb:/usr/local/var/lib/couchdb --name dondedb --env COUCHDB_USER=donde --env COUCHDB_PASSWORD=dondepwd couchdb
        docker run --detach --volume /usr/local/var/lib/couchdb:/usr/local/var/lib/couchdb --name ${name} couchdb
    fi
fi
