package main

import (
	"flag"
	"github.com/hokiegeek/donde-estas-daemon"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "", 0)
	portPtr := flag.Int("port", 8585, "Specify the port to use")
	flag.Parse()

	logger.Printf("Serving on port %d\n", *portPtr)

	dondeestas.New(logger, *portPtr)
}
