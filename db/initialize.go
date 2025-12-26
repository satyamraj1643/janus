package db


import (
	"context"
	"log"
	"os"
	"time"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Init(){

	dsn := os.Getenv("DB_URL")

	if dsn == ""{
		log.Fatal("DB URL not set.")
	}

	config, err := pgxpool.ParseConfig(dsn)

	if err != nil {
		log.Fatal(err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour

	Pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}

	if err := Pool.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}

	log.Println("Postgre Instance connected")

}