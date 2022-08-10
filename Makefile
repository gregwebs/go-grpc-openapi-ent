proto_source_files := $(wildcard ./apis/todo/v1/*.proto apis/*.yaml)

gen/proto/go/todo/v1/*.go: $(proto_source_files)
	just buf generate

ent_source_files := $(wildcard ./ent/schema/*.go ent/generate.go)

ent/ent.go: $(ent_source_files)
	go generate ./ent
