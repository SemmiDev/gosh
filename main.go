package main

import (
	"gosh/internal"
	"gosh/internal/postgres"
	"gosh/internal/program"
	"log"
)

func main() {
	pgxPoolConn, err := postgres.NewPgxPoolConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer pgxPoolConn.Close()

	programDataStore := program.NewProgramDataStore(pgxPoolConn)
	programSearchService := program.NewProgramSearchService(programDataStore)
	programHttpHandler := program.NewHttpHandler(programDataStore, programSearchService)
	httpServer := internal.NewHttpServer(
		programHttpHandler,
		programSearchService,
	)

	httpServer.SetupHttpMiddleware()
	httpServer.SetupWebSocketMiddleware()
	httpServer.SetupHttpRoutes()
	httpServer.SetupWebSocketRoutes()
	httpServer.RunWithGracefulShutdown()
}
