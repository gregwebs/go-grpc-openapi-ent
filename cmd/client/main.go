package main

import (
	"context"
	"fmt"
	"log"

	// This import path is based on the name declaration in the go.mod,
	// and the gen/proto/go output location in the buf.gen.yaml.
	apiv1 "github.com/gregwebs/go-grpc-openapi-ent/gen/proto/go/todo/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	connectTo := "127.0.0.1:8080"
	conn, err := grpc.Dial(
		connectTo,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to todoService on %s: %w", connectTo, err)
	}
	log.Println("Connected to", connectTo)

	svc := apiv1.NewTodoServiceClient(conn)
	ctx := context.Background()
	rspEmps, err := svc.ListTaskUsers(ctx, &apiv1.ListTaskUsersRequest{
		Id: 2,
	})
	if err != nil {
		return fmt.Errorf("failed listing task user")
	}
	for _, emp := range rspEmps.Users {
		if _, err := svc.GetUser(ctx, &apiv1.GetUserRequest{
			Id:    0,
			Email: emp.Email,
		}); err != nil {
			return fmt.Errorf("failed to GetUser: %w", err)
		}
		log.Printf("Successfully Got User %s", emp.Email)
	}

	return nil
}
