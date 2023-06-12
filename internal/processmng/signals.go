package processmng

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func OnInterrupt(cbs ...func()) {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	// waiting for receiving signals
	<-sigs
	log.Default().Println("\ninterrupt signal received")
	for _, cb := range cbs {
		cb()
	}
}
