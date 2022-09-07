#!/usr/bin/env just --justfile

default:
	just -l

build: generate
	cd server && go build
	@echo "server built"
	cd cmd/client && go build
	@echo "client built"
	cd cmd/migrate && go build
	@echo "migrate built"
	cd cmd/seed && go build
	@echo "seed built"
	@just lint

lint:
	golangci-lint run
	staticcheck ./...

release:
	just buf breaking --against ".git#branch=main,subdir=."

setup: setup-db setup-go generate

generate:
	make gen/proto/go/todo/v1/*.go >/dev/null
	make ent/ent.go >/dev/null

run-server: generate
	cd server && go build
	DB_NAME={{db_name}} DB_USER={{db_user}} ./server/server

run-client: generate
	cd cmd/client && go build
	./cmd/client/client

run-migrate *args='': generate
	cd cmd/migrate && go build
	DB_NAME={{db_name}} DB_USER={{db_user}} ./cmd/migrate/migrate {{args}}

run-seed *args='': generate
	cd cmd/seed && go build
	DB_NAME={{db_name}} DB_USER={{db_user}} ./cmd/seed/seed {{args}}

buf *args='':
	cd apis && PATH="$PATH:$(go env GOPATH)/bin" buf {{args}}

ent *args='':
	go run entgo.io/ent/cmd/ent {{args}}

db_name := 'todo'
db_user := 'todo'
setup-db:
	echo "SELECT 'CREATE DATABASE {{db_name}}' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '{{db_name}}')\gexec" | psql postgres # 2>/dev/null
	echo "SELECT 'CREATE USER {{db_user}}' WHERE NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = '{{db_user}}')\gexec" | psql postgres # 2>/dev/null

setup-go:
	@# Doesn't seem to work from the gen directory
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install entgo.io/ent/cmd/ent@latest
	cd gen && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	cd gen && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	cd gen && go install \
	  github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
	  github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	cd gen && go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
