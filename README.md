# Golang Microservice

The objective is a blue-print for a production-ready microservice written in Go, designed to be deployed in Kubernetes.
It connects to a Kafka topic and uses a PostgreSQL database.

## Project Structure

- `api/`: Contains the OpenAPI specification file.
- `cmd/server/`: Contains the main entry point of the microservice.
- `pkg/`: Contains the main packages of the microservice.
  - `config/`: Handles configuration management.
  - `db/`: Defines the PostgreSQL database connection and migrations.
  - `kafka/`: Implements the Kafka consumer functionality.
  - `handler/`: Defines the HTTP handlers for the microservice endpoints.
  - `model/`: Contains the data models and structs.
  - `repository/`: Implements the data access layer.
  - `service/`: Contains the business logic and service layer.
- `internal/util/`: Contains utility packages used internally.
- `build/`: Contains the Dockerfile for building the microservice container image.
- `deploy/kubernetes/`: Contains Kubernetes deployment files.
- `scripts/`: Contains database migration scripts.
- `test/integration/`: Contains integration tests for the microservice.

## Getting Started

1. Install the necessary dependencies:
   - Go
   - Swagger code generator (`go-swagger`)
   - PostgreSQL
   - Kafka

2. Generate the server code from the OpenAPI specification: 
   ``` swagger generate server -f api/openapi.yaml -A golang-kafka-kubernetes-microservice -t pkg/handler/generated ```
3. Update the configuration in the .env file.
4. Build and run the microservice: ``` go build -o golang-kafka-kubernetes-microservice cmd/server/main.go ./golang-kafka-kubernetes-microservice ```
5. Access the API endpoints using the generated Swagger UI or by sending requests to the appropriate URLs.
   
## Database Migrations
The microservice uses database migrations. Migration files are located in pkg/db/migration/. To create a new migration, add a new SQL file with the desired changes.

## Kafka Topic Creation
The microservice automatically creates the Kafka topic if it doesn't exist. The topic configuration can be modified in pkg/kafka/topic.go.

## Deployment
The microservice is designed to be deployed in Kubernetes. The necessary deployment files are located in deploy/kubernetes/. Adjust the configuration as needed for your Kubernetes environment.

## Testing
Integration tests are located in test/integration/. Run the tests using the following command:
```
go test ./test/integration/...
```