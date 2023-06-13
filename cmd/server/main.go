package main

import (
	"log"

	"github.com/ARTM2000/archive1/internal/archive"
)

func main() {
	log.Default().Println("[ Archive1 Server ]")
	archive.CmdExecute()
}
