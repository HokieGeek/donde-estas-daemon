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

	db, err := dondeestas.NewDbClient(dondeestas.CouchDB)
	if err != nil {
		panic(err)
	}

	dondeestas.New(logger, *httpPortPtr, db)
}
