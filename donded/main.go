package main

import (
	"flag"
	"github.com/hokiegeek/donde-estas-daemon"
	"log"
	"os"
)

func main() {
	httpPortPtr := flag.Int("port", 8585, "Specify the port to use")
	flag.Parse()

	logger := log.New(os.Stdout, "", 0)
	logger.Printf("Serving on port %d\n", *httpPortPtr)

	params := dondeestas.DbClientParams{dondeestas.CouchDB, "donde", "db", 5984}

	db, err := dondeestas.NewDbClient(params)
	if err != nil {
		panic(err)
	}

	dondeestas.New(logger, *httpPortPtr, db)
}
