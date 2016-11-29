#!/bin/sh

docker build -t hokiegeek/godonde .&&docker run -it --rm -p 8585:8585 --link dondedb:db hokiegeek/godonde
