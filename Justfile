default:
    just --list

build:
    go build -o bin/server.exe cmd/server/main.go

run:
    ./bin/server.exe