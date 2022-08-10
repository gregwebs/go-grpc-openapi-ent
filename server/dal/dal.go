// dal = Data Access Layer
package dal

import (
	"context"
	"fmt"
	"log"

	ent "github.com/gregwebs/go-grpc-openapi-ent/ent"
	"github.com/gregwebs/go-grpc-openapi-ent/ent/task"
	"github.com/gregwebs/go-grpc-openapi-ent/ent/user"
	_ "github.com/lib/pq"
)

func UserByEmail(ctx context.Context, db *ent.Client, email string) (*ent.User, error) {
	return db.User.Query().Where(
		user.Email(email),
	).First(ctx)
}

func TaskUsers(ctx context.Context, db *ent.Client, taskID int) (*ent.Task, error) {
	return db.Task.Query().WithUsers().Where(
		task.IDEQ(taskID),
	).First(ctx)
}

func WithTx(ctx context.Context, client *ent.Client, fn func(tx *ent.Tx) error) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if v := recover(); v != nil {
			if errRoll := tx.Rollback(); errRoll != nil {
				log.Printf("error rolling back: %+v", errRoll)
			}
			panic(v)
		}
	}()
	if err := fn(tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = fmt.Errorf("rolling back transaction: %w", rerr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}
