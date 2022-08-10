package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	// "entgo.io/ent/dialect/sql/schema"

	migrate "github.com/gregwebs/go-grpc-openapi-ent/ent/migrate"
	dal "github.com/gregwebs/go-grpc-openapi-ent/server/dal"

	"ariga.io/atlas/sql/sqltool"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	// _ "ariga.io/atlas/sql/postgres"
	// _ "github.com/jackc/pgx/v4"
)

var (
	// command-line options:
	// gRPC server endpoint
	flagPrint = flag.Bool("print", false, "print migration sql")
	flagDSN   = flag.Bool("dsn", false, "show the DSN connection string")
)

func main() {
	flag.Parse()
	ctx := context.Background()
	if *flagPrint {
		printMigration()
		return
	}

	dsn := dal.ConnectDSN("postgres", dal.ConnectConf{
		User:     "",
		DB:       "",
		Password: "",
	})
	fmt.Println(dsn)
	if *flagDSN {
		return
	}

	// Create a local migration directory able to understand golang-migrate migration files for replay.
	dir, err := sqltool.NewGolangMigrateDir("ent/migrations")
	if err != nil {
		log.Fatalf("failed creating atlas migration directory: %v", err)
	}
	// Write migration diff.
	opts := []schema.MigrateOption{
		schema.WithDir(dir),                          // provide migration directory
		schema.WithMigrationMode(schema.ModeInspect), // provide migration mode
		schema.WithDialect(dialect.Postgres),         // Ent dialect to use
	}

	err = migrate.NamedDiff(ctx, dsn, "todo", opts...)
	if err != nil {
		log.Fatalf("failed generating migration file: %v", err)
	}
}

func printMigration() {
	db, err := dal.Connect(dal.ConnectConf{
		User:     "",
		DB:       "",
		Password: "",
	})
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	ctx := context.Background()
	/*
		f, err := os.Create("migrate.sql")
		if err != nil {
			log.Fatalf("create migrate file: %v", err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				fmt.Println(err)
			}
		}()
	*/
	opts := []schema.MigrateOption{
		schema.WithDropColumn(true),
		schema.WithDropIndex(true),
	}
	if err := db.Schema.WriteTo(ctx, os.Stderr, opts...); err != nil {
		log.Fatalf("failed printing schema changes: %v", err)
	}

	if err != nil {
		log.Fatalf("migration error %+v", err)
	}
}
