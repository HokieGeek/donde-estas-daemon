package main

import (
	"flag"
	"github.com/hokiegeek/donde-estas-daemon"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	httpPortPtr := flag.Int("port", 8080, "Specify the port to use")
	databaseUrlPtr := flag.String("dburl", "db:5984", "The hostname[:port] of the database")
	flag.Parse()

	sepPos := strings.LastIndex(*databaseUrlPtr, ":")
	dbHost := (*databaseUrlPtr)[:sepPos]
	dbPort, _ := strconv.Atoi((*databaseUrlPtr)[sepPos+1:])

	logger := log.New(os.Stdout, "", 0)
	logger.Printf("Connecting to %s on port %d\n", dbHost, dbPort)
	logger.Printf("Serving on port %d\n", *httpPortPtr)

	params := dondeestas.DbClientParams{dondeestas.CouchDB, "donde", dbHost, dbPort}

	db, err := dondeestas.NewDbClient(params)
	if err != nil {
		panic(err)
	}

	logger.Fatal(dondeestas.ListenAndServe(logger, *httpPortPtr, db))
}
