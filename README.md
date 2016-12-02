# ¿Dónde Estás? Daemon [![Build Status](https://travis-ci.org/HokieGeek/donde-estas-daemon.svg?branch=master)](https://travis-ci.org/HokieGeek/donde-estas-daemon) [![Coverage](http://gocover.io/_badge/github.com/HokieGeek/donde-estas-daemon)](http://gocover.io/github.com/HokieGeek/donde-estas-daemon) [![GoDoc](http://godoc.org/github.com/HokieGeek/donde-estas-daemon?status.png)](http://godoc.org/github.com/HokieGeek/donde-estas-daemon)
The server side to the [¿Dónde Estás?](https://github.com/HokieGeek/DondeEstas) android app

##### Suggested usage
###### Starting server
```sh
docker run -d --name couchdb couchdb
docker run -d -p 8080:8080 --link couchdb:db hokiegeek/donde-estas-daemon
```

###### Adding an entry
```sh
curl localhost:8080/update --data \
'{
    "id" : "42",
    "name" : "foobar",
    "position": {
        "tov": "2009-11-10T23:00:00Z",
        "latitude": 0.1,
        "longitude": 0.2,
        "elevation": 0.3
    },
    "visible" : true,
    "whitelist" : ["99"],
    "following" : []
}'

curl localhost:8080/update --data \
'{
    "id" : "99",
    "name" : "raboof",
    "position": {
        "tov": "2009-11-10T22:00:00Z",
        "latitude": 0.4,
        "longitude": 0.5,
        "elevation": 0.6
    },
    "visible" : true,
    "whitelist" : ["42"],
    "following" : ["42"]
}'
```

###### Retrieving entries
```sh
curl localhost:8080/person --data '{"ids":["42","99"]}'
```
