#!/bin/bash

DATABASE_URL="postgresql://user:password@localhost:5432/mydb?sslmode=disable"

MIGRATION_DIR="pkg/db/migration"

migrate -path $MIGRATION_DIR -database $DATABASE_URL up

migrate -path $MIGRATION_DIR -database $DATABASE_URL version