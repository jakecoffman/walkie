package main

import (
	"github.com/jakecoffman/walkie"
	"github.com/jakecoffman/walkie/eng"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile)
	eng.Run(&walkie.Game{})
}
