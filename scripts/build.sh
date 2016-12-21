#!/bin/sh

docker build -t hokiegeek/donde-estas-daemon . && \
docker run -it --rm --publish 8585:8080 --link couchdb:db hokiegeek/donde-estas-daemon
