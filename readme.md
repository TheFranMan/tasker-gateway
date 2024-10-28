# Tasker Gateway

This application serves as an API gateway to the Tasker ecosystem. End users can utilize various APIs to log requests in the tasker backend and then query the status of the request.

## Dependencies
   * **Mysql**: Persist the requests.
   * **Redis**: Store the status for a token.
   * **Prometheus**: Record metrics.
   * **Grafana**: Visualize the gathered Prometeus metrics.

## Run

In the project root, first create a `.env` file with the following env vars below. Then run `docker compose up -d`.

| Name          | Required           | Default | Description                                |
| ------------- | ------------------ | ------- | ------------------------------------------ |
| PORT          | :white_check_mark: | n/a     | The port the server runs on.               |
| ENV           | :x:                | n/a     | The enviroment the app is running in.      |
| DB_USER       | :white_check_mark: | n/a     | The local MySQL user.                      |
| DB_PASS       | :white_check_mark: | n/a     | The local MySQL password.                  |
| DB_HOST       | :white_check_mark: | n/a     | The local MySQL host.                      |
| DB_PORT       | :white_check_mark: | n/a     | The local MySQL port.                      |
| DB_NAME       | :white_check_mark: | n/a     | The local MySQL database name.             |
| AUTH_TOKENS   | :white_check_mark: | n/a     | List of API auth tokens sepearted by a "," |
| REDIS_ADDR    | :white_check_mark: | n/a     | The local Redis address                    |
| REDIS_KEY_TTL | :x:                | 30s     | The TTL for the status Redis key           |



To run the project sepearte from Docker, and so avaoiding the Docker build step, from the project root follow the below steps:
   * Comment out the `gateway` block from the `docker-compose.yml` file.
   * Run `docker compose up -p` to start its dependencies.
   * Run `go get && go build && ./gateway`

## Testing
Tests can be run from the root dir using: `go test -v ./... --tags=integration`. The `integration` tag runs the integration tests using Docker. Remove this if you only want to run unit tests.