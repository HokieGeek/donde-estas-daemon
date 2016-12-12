#!/bin/sh

docker build -t hokiegeek/donde-estas-daemon . && \
docker run -it --rm -p 8585:8080 --link dondedb:db hokiegeek/donde-estas-daemon
