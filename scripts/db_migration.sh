#!/bin/bash

# Set the database connection URL
DATABASE_URL="postgresql://user:password@localhost:5432/mydb?sslmode=disable"

# Set the path to the migration files
MIGRATION_DIR="pkg/db/migration"

# Run the database migrations
migrate -path $MIGRATION_DIR -database $DATABASE_URL up

# Check the status of the migrations
migrate -path $MIGRATION_DIR -database $DATABASE_URL version