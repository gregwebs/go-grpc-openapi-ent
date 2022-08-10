# Overview

This is a template for creating a type-safe Go API. It relies heavily on code generation. In languages without generics (just added to Go) and macros, type-safety is achieved by a combination of reflection and code generation. Code generation works well for Go because it compiles and links very quickly.

The API is created with GPRC and the GRPC HTTP Gateway. This extends the type-safety to clients: client code is generated for GRPC or HTTP.

The API is described in a [.proto file](./apis/todo/v1/todo.proto). The code generation process generates OpenAPI, and the HTTP Gateway allows for JSON rather than be restricted to only GRPC. The JSON API is documented at `http://localhost:8081/openapi/`. Put in `/openapi.json` in the box at the top of the page.

For database access this project uses the Ent framework from Meta. It is based around describing a database schema in Golang: this is done in the folder [./schema](./schema). In the past I have also used GORM + go-queryset to good effect. The biggest initial improvement that Ent provided for me is versioned SQL migration that can be inspected, altered, and applied.

Generated code is under [./gen](./gen).
The go code for this project is under the directories [./server](./server) and [./cmd](./cmd).

There is an example client program under [./cmd/client](./cmd/client).

This project uses the `just` command runner to automate command execution. Often this is done in projects with a `Makefile`, but `make` is not designed for this task: it is designed for tracking build dependencies. There is a [Makefile](./Makefile) for this project that is used to track build dependencies to avoid re-running time-consuming code generation steps: this allows code generation to be ran automatically instead of being a manual process.

# TLS and auth

TODO: Currently this template does not use TLS or any auth: you must add this yourself.

# Modifying

To make this code your own, fork/clone/copy the repo and then update the repo name (the below uses [sad](https://github.com/ms-jpq/sad) for find and replace).

	find ./ -type f -name '*.go' | sad gregwebs/go-grpc-openapi-ent ghuser/ghrepo

Then change the LICENSE file to match your project needs (this code is committed to the public domain).

Then you can alter the [.proto api description file](./apis/todo/v1/todo.proto). You can also rename it, and its folder, and alter some of the buf*.yaml files accordingly.

Alter the database by changing the files in the [./schema](./schema) folder.

# Installation


Install Go if you don't have it already, e.g.

	brew install go

Install dev dependencies that are listed in the [Brewfile](./Brewfile)

	brew bundle

This will tools required for local development.

* just - a command runner (rather than using Make, which is a build tool)
* buf - a cli to generate protobuf and gRPC code from .proto definitions

But you can also install these packages a different way
Please do still install `just` with `brew install just` or

	curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | bash -s -- --to ./
	chmod +x ./just
	mv just <where you put executables>

You also need to install and run your database of choice.
For example:

	brew install postgresql
	$HOMEBREW_PREFIX/bin/postgres -D $HOMEBREW_PREFIX/var/postgres


# Install go tools

	just setup-go


# Database migration

First setup your database. There is a small script for doing this with Postgres in the Justfile via

	just setup-db

You can then run the migrations with:

	migrate -source file://ent/migrations -database <DSN> up

`migrate` is a go tool that should be installed by `just setup`.
The DSN connection string can be printed with:

	just run-migrate -dsn

New versioned migrations are generated to the gen/migrations folder with:

	mkdir ent/migrations
	just run-migrate


# Database seeding

	just run-seed


# Running the server

	just run-server


# Using the api

There is an example GRPC client interaction that can be run:

	just run-client

You can also use the HTTP API on port 8081. See the above instructions for opending the OpenAPI documentation. Requests from the OpenAPI documentation itself doesn't work because it assumes `https://`, but it will also show instructions for how to use curl, and you can change that to `http://`

# Building and linting

The above commands will re-build what they need. To build all commands and run the linter, run `just build`. To run the linter, run `just lint`. 
