#!/bin/bash

set -euo pipefail

set -o allexport
source .env set
set +o allexport

docker run --rm -v "./migration" --network bot-reminder migrate/migrate -source=file:///migration -database "postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_DB?sslmode=disable"  $@
