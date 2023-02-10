package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
	"runtime"

	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	connConfig, err := pgx.ParseConfig("postgres://root:secret@localhost/gosh")
	if err != nil {
		log.Fatalln("Unable to parse config:", err)
		return
	}

	pgxPoolConfig := pgxpool.Config{
		ConnConfig: connConfig,
		MaxConns:   int32(runtime.NumCPU() * 5),
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), &pgxPoolConfig)
	if err != nil {
		log.Fatalln("Unable to connect to database:", err)
		return
	}
	defer pool.Close()

	fmt.Println("Successfully connected to database")
}
