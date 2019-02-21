package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	signal.Notify(gracefulStop, syscall.SIGABRT)

	server := Server{}
	go func(server *Server) {
		sig := <-gracefulStop
		fmt.Printf("caught sig: %+v", sig)
		server.database.db.Close()
		fmt.Println("Closing database.")
		os.Exit(0)
	}(&server)
	server.init()
}
