# Arvut DB API (gxydb-api)
Backend API for the BB galaxy system.

## Setup

1. Clone the repository
2. Create a `.env` file with the following environment variables:
   ```
   SECRET=12345678901234567890123456789012
   MQTT_BROKER_URL=localhost:1883
   ```
   Note: The SECRET must be exactly 32 bytes for AES-256 encryption.

3. Run database migrations:
   ```
   go run cmd/migrate.go
   ```

4. Generate models using SQLBoiler:
   ```
   sqlboiler psql
   ```

## Building

```
go build
```

## Running Tests

To run tests that don't require external services:
```
./run_tests.sh
```

To run all tests (requires MQTT broker and Janus Gateway):
```
go test ./...
```

## External Dependencies

Some tests require external services:
- MQTT broker running on localhost:1883
- Janus Gateway with Admin API running on localhost:7088

If these services are not available, some tests will fail.


### Development Environment
In local dev env we use docker-compose to have our dependencies setup for us.

Fire up all services
```shell script
docker-compose up -d
```


CLI access to local dev db
```shell script
docker-compose exec db psql -U user -d galaxy
```

#### Try out the API
To try out authenticated endpoints on your local environment set the `SKIP_AUTH=true` environment variable.
This will allow every request to any endpoint.

However, on production, you must send a valid `Authorization` header.
You can get one by inspecting network traffic in browser after login in. 
Copy paste the value of the `Authorization` header and use that.


### DB Migrations

DB Migrations are managed by [migrate](https://github.com/golang-migrate/migrate)

On mac install with brew. For other platforms see the project homepage.
```shell script
brew install golang-migrate
```

Create a new migration by
```shell script
migrate create -ext .sql -dir migrations -format 20060102150405 <migration_name_goes_here>
```

Run migrations
```shell script
migrate -database "postgres://user:password@localhost/galaxy?sslmode=disable" -path migrations up
```

Regenerate models

```shell script
sqlboiler psql
```

Download instructions for sqlboiler can be found [here](https://github.com/volatiletech/sqlboiler#download).

