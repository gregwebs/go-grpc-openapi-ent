package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	// This import path is based on the name declaration in the go.mod,
	// and the gen/proto/go output location in the buf.gen.yaml.
	ent "github.com/gregwebs/go-grpc-openapi-ent/ent"
	apiv1 "github.com/gregwebs/go-grpc-openapi-ent/gen/proto/go/todo/v1"
	gw "github.com/gregwebs/go-grpc-openapi-ent/gen/proto/go/todo/v1"
	api "github.com/gregwebs/go-grpc-openapi-ent/server/api"
	dal "github.com/gregwebs/go-grpc-openapi-ent/server/dal"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	// command-line options:
	// gRPC server endpoint
	grpcEndpoint = flag.String("grpc-server-endpoint", "localhost:8080", "gRPC server endpoint")
)

func main() {
	flag.Parse()
	db, err := dal.Connect(dal.ConnectConf{
		User:     "",
		DB:       "",
		Password: "",
	})
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	go func() {
		for i := 0; i < 10; i++ {
			if err := runHTTP(); err != nil {
				log.Printf("%+v", err)
			}
			time.Sleep(time.Millisecond * 100)
		}
		log.Fatal("HTTP Gateway shutdown")
	}()

	if err := runGRPC(db); err != nil {
		log.Fatal(err)
	}
	log.Fatal("GRPC shutdown")
}

func runHTTP() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux()

	// Show openapi documentation.
	// Go to /openapi/
	// In the box at the top of the page, put in /openapi.json
	fs := http.FileServer(http.Dir("./apis/ui"))
	if err := mux.HandlePath("GET", prefixOpenApi+".json", handleOpenapiDescription); err != nil {
		log.Fatal("add to mux openapi json", err)
	}
	if err := mux.HandlePath("GET", prefixOpenApi+"/*", handleHandler(http.StripPrefix(prefixOpenApi, fs))); err != nil {
		log.Fatal("add to mux openapi", err)
	}

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := gw.RegisterTodoServiceHandlerFromEndpoint(ctx, mux, *grpcEndpoint, opts)
	if err != nil {
		return err
	}
	listenOn := ":8081"
	log.Println("Listening for HTTP on", listenOn)
	return http.ListenAndServe(listenOn, mux)
}

func runGRPC(db *ent.Client) error {
	listenOn := *grpcEndpoint
	listener, err := net.Listen("tcp", listenOn)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", listenOn, err)
	}

	server := grpc.NewServer()
	apiv1.RegisterTodoServiceServer(server, api.NewtodoService(
		db,
	))
	log.Println("Listening on", listenOn)
	if err := server.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	}

	return nil
}

func handleHandler(h http.Handler) func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	return func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		h.ServeHTTP(w, r)
	}
}

func handleOpenapiDescription(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	content, err := os.ReadFile("gen/openapiv2/todo/v1/todo.swagger.json")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	if _, err = w.Write(content); err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
}

const prefixOpenApi = "/openapi"
