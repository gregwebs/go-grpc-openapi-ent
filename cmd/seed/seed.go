package main

import (
	"context"
	"fmt"
	"log"

	ent "github.com/gregwebs/go-grpc-openapi-ent/ent"
	dal "github.com/gregwebs/go-grpc-openapi-ent/server/dal"
)

func main() {
	if err := seed(); err != nil {
		log.Fatalf("%+v", err)
	}
}

func seed() error {
	ctx := context.Background()

	db, err := dal.Connect(dal.ConnectConf{
		User:     "",
		DB:       "",
		Password: "",
	})
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	fmt.Println("seeding the DB")
	if err := seedDB(ctx, db); err != nil {
		return err
	}
	return nil
}

func seedDB(ctx context.Context, db *ent.Client) error {
	return dal.WithTx(ctx, db, func(tx *ent.Tx) error {
		inflationTask := db.Task.Create().SetDescription(
			"Make TODO LisT",
		).SaveX(ctx)
		reserveTask := db.Task.Create().SetDescription(
			"Do first item on TODO LIST",
		).SaveX(ctx)
		btcTask := db.Task.Create().SetPrivate(true).SetDescription(
			"Do second item on TODO List",
		).SaveX(ctx)
		digitalTask := db.Task.Create().SetDescription(
			"Do third item on TODO list",
		).SaveX(ctx)

		bulkEmps := []*ent.UserCreate{
			db.User.Create().SetName("alice").SetEmail("alice@example.com").
				AddTasks(inflationTask),
			db.User.Create().SetName("bob").SetEmail("bob@example.com").
				AddTasks(reserveTask, btcTask),
			db.User.Create().SetName("charlie").SetEmail("charlie@example.com").
				AddTasks(btcTask, digitalTask),
		}
		_, err := db.User.CreateBulk(bulkEmps...).Save(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}
