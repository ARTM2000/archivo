package main

import (
	"log"

	"github.com/ARTM2000/archivo/internal/archive"
)

func main() {
	log.Default().Println("[ Archivo Server ]")
	archive.CmdExecute()
}
