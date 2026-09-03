#!/usr/bin/env bash

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

"$DIR/my_program" "$@"

trap "trap - SIGTERM && kill -- -$$" SIGINT SIGTERM EXIT

go run "$DIR/cmd/raft-node/main.go" --id 1 --grpc-port 8081 --http-port 8091 --peers=localhost:8082,localhost:8083 "$@" &
go run "$DIR/cmd/raft-node/main.go" --id 2 --grpc-port 8082 --http-port 8092 --peers=localhost:8081,localhost:8083 "$@" &
go run "$DIR/cmd/raft-node/main.go" --id 3 --grpc-port 8083 --http-port 8093 --peers=localhost:8081,localhost:8082 "$@" &

wait
