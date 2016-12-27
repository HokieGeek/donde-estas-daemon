#!/bin/bash

here=$(cd ${0%/*}; pwd)

host="127.0.0.1" #$(hostname --fqdn)
port=8585

getRandomString() {
    head /dev/urandom | tr -dc A-Za-z0-9 | head -c $((RANDOM%30))
}

verifyPerson() {
    id=$1; shift
    while (( $# > 0 )); do
        key=${1%%:*}
        val=${1##*:}
        case ${key} in
        lat) lat=${val} ;;
        lon) lon=${val} ;;
        follow) following="${val//,/\",\"}" ;;
        white) whitelist="${val//,/\",\"}" ;;
        esac
        shift
    done

    person=$(${here}/person.sh get ${host}:${port} ${id})
    # echo "PERSON: ${person}"

    foundId="$(echo ${person} | sed 's/.*"id":"\([^"]*\)".*/\1/')"
    # set -x
    [ "${id}" != "${foundId}" ] && return 1

    # delta=.000005

    # foundLat=$(echo ${person} | sed 's/.*"latitude":\([^,]*\),.*/\1/')
    # (( foundDiff = ${lat} - ${foundLat} ))
    # (( ${foundDiff} < ${delta} )) && return 2
    # (( ${lat} - ${foundLat} < ${delta} )) && return 2
    # [[ ${lat} != ${foundLat}* ]] && return 2

    # foundLon=$(echo ${person} | sed 's/.*"longitude":\([^,]*\),.*/\1/')
    # [ "${lon}" != "${foundLon}" ] && return 3

    foundFollowing=$(echo ${person} | sed 's/.*"following":\[\([^]]*\).*/\1/;s/"//g')
    [ "${following}" != "${foundFollowing}" ] && return 4

    foundWhitelist=$(echo ${person} | sed 's/.*"whitelist":\[\([^]]*\).*/\1/;s/"//g')
    [ "${whitelist}" != "${foundWhitelist}" ] && return 5
}

cleanOnExit() {
    set +x
    docker stop ${servername} && docker rm ${servername}
    docker stop ${dbname} && docker rm ${dbname}
    sudo rm -rf ${dbDir}
}
trap cleanOnExit EXIT

ids+=($(getRandomString))
ids+=($(getRandomString))

lat="39.18"$(shuf -i 0-9999 -n 1)
lon="-77.27"$(shuf -i 0-9999 -n 1)

updateOne="${ids[0]} lat:${lat} lon:${lon} white:${ids[1]}"
updateTwo="${ids[1]} lat:${lat} lon:${lon} follow:${ids[0]} white:${ids[0]}"

dbname="couchdb_$$"
servername="donded_$$"

## Create the database
dbDir=/usr/local/var/lib/couchdb
sudo mkdir ${dbDir}
docker run --detach --volume ${dbDir}:${dbDir} --name ${dbname} couchdb:1.6

## Create the server
docker build --tag hokiegeek/donde-estas-daemon ${here}/..
docker run --detach --name ${servername} --publish ${port}:8080 --link ${dbname}:db hokiegeek/donde-estas-daemon

sleep 10s

docker logs -f ${dbname} >&3 2>&1 &
docker logs -f ${servername} >&4 2>&1 &

## Add some "real" data
${here}/person.sh update ${host}:${port} ${updateOne}
${here}/person.sh update ${host}:${port} ${updateTwo}

## Query that data
verifyPerson ${updateOne} || exit $?
verifyPerson ${updateTwo} || exit $?
