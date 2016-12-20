package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/hokiegeek/donde-estas-daemon"
)

func main() {
	httpPortPtr := flag.Int("port", 8080, "Specify the port to use for the HTTP server")
	databaseURLPtr := flag.String("dburl", "db:5984", "The hostname:port of the database server")
	databaseNamePtr := flag.String("dbname", "donde", "The name of the database")
	flag.Parse()

	sepPos := strings.LastIndex(*databaseURLPtr, ":")
	dbHost := (*databaseURLPtr)[:sepPos]
	dbPort, _ := strconv.Atoi((*databaseURLPtr)[sepPos+1:])

	logger := log.New(os.Stdout, "", 0)
	logger.Printf(":: Connecting to database %s at %s on port %d\n", dbHost, dbPort)
	logger.Printf(":: Serving on port %d\n", *httpPortPtr)

	params := dondeestas.DbClientParams{dondeestas.CouchDB, *databaseNamePtr, dbHost, dbPort}

	db, err := dondeestas.NewDbClient(params)
	if err != nil {
		panic(err)
	}

	logger.Fatal(dondeestas.ListenAndServe(logger, *httpPortPtr, db))
}
