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
	dbPort, _ := strconv.ParseUint((*databaseURLPtr)[sepPos+1:], 10, 16)

	logger := log.New(os.Stdout, "", 0)
	logger.Printf(":: Connecting to database %s at %s on port %d\n", *databaseNamePtr, dbHost, dbPort)
	logger.Printf(":: Serving on port %d\n", *httpPortPtr)

	db, err := dondeestas.NewDbClient(dondeestas.DbClientParams{
		DbType:   dondeestas.CouchDB,
		DbName:   *databaseNamePtr,
		Hostname: dbHost,
		Port:     uint16(dbPort),
	})
	if err != nil {
		panic(err)
	}

	logger.Fatal(dondeestas.ListenAndServe(logger, *httpPortPtr, db))
}
