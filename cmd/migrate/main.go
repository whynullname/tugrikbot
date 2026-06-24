package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/whynullname/tugrikbot/migrations"
)

func main() {
	dsn := os.Getenv("TUGRIK_DB_DSN")

	if dsn == "" {
		log.Fatal("empty dsn in env")
	}

	if len(os.Args) < 2 {
		log.Fatal("команда не указана: up|down")
	}

	command := os.Args[1]
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("open connect error: %v", err)
	}
	defer db.Close()

	err = db.PingContext(context.Background())
	if err != nil {
		log.Fatalf("ping db error: %v", err)
	}

	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("goose error: %v", err)
	}

	err = goose.RunContext(context.Background(), command, db, ".")
	if err != nil {
		log.Fatalf("goose error: %v", err)
		os.Exit(1)
	}
}
