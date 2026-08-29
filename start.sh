#!/usr/bin/env bash

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

"$DIR/my_program" "$@"

trap "trap - SIGTERM && kill -- -$$" SIGINT SIGTERM EXIT

go run "$DIR/cmd/raft-node/main.go" --id 1 --port 8081 --peers=localhost:8082,localhost:8083 "$@" &
go run "$DIR/cmd/raft-node/main.go" --id 2 --port 8082 --peers=localhost:8081,localhost:8083 "$@" &
go run "$DIR/cmd/raft-node/main.go" --id 3 --port 8083 --peers=localhost:8081,localhost:8082 "$@" &

wait
