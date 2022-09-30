#!/bin/bash
set -e

DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"


(docker stop postgres && docker rm postgres) || true

sudo rm -rf $DIR/data

docker run \
    -d \
    -p 5432:5432 \
    --name postgres \
    -e POSTGRES_PASSWORD=P@ssw0rd \
    -e PGDATA=/var/lib/postgresql/data \
    -v $(pwd)/data:/var/lib/postgresql/data \
    -v $(pwd)/init:/docker-entrypoint-initdb.d \
    postgres:14.5
