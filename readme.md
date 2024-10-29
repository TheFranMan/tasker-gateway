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

## Authenication
The `/api` routes are authenicated by a `Authorization` header with a valid token. These tokens are specified in the env vars. **NOTE**: This is only a placeholder implementation to show the possible statuses the authenication middleware can provide. In reality, a more robust authenication method like OAuth2, AzureAD, or another third party system would be used.

## Load Testing
The [locust](https://locust.io/) application is a Python based tool used for load testing. Follow the [installation steps](https://locust.io/#install) to install locust, then from the root dir run `locust`.
This will bring up a web UI located at [http://0.0.0.0:8089](http://0.0.0.0:8089) where you can specify the parameters of the test. The tests are located in the file `locustfile.py` in the root dir.

## Monitoring
Monitoring is supplied via Prometheus and using a Grafana dashboard at [http://localhost:3021/](http://localhost:3021/). The username and password are both `admin`. The dashboard can be imported using the file `grafana-dashboard.json` in the root dir.

## Routes
   * `GET /heartbeat`

     #### Description:

     Simple endpoint used to return a 200 response to check for application health.

     #### Responses:
     * 200

   * `GET /metrics`

     #### Description:

     Used by the Prometheus scrapper to gather site statistics

     #### Responses:
     * 200

        Body:
        ```
       # HELP go_gc_duration_seconds A summary of the wall-time pause (stop-the-world) duration in garbage collection cycles.
       # TYPE go_gc_duration_seconds summary
       go_gc_duration_seconds{quantile="0"} 0.000521397
       go_gc_duration_seconds{quantile="0.25"} 0.000521397
       go_gc_duration_seconds{quantile="0.5"} 0.000521397
       go_gc_duration_seconds{quantile="0.75"} 0.000521397
       go_gc_duration_seconds{quantile="1"} 0.000521397
       ...
       ```

   * `DELETE /api/user`

      #### Description:
      Add a delete request on the tasker backend. A token is returned in the body that can be used to poll for the status of the request.

      #### Responses:

         * 201

            Body:
            ```json
            {
               "token": "abc"
            }
            ```

         * 400

            Body:
            ```
            The reason in plain text
            ```

         * 401

         * 415

         * 500

            Body:
            ```
            The reason in plain text
            ```

   * `GET /api/poll/{token}`

      #### Description:
      Poll for the status of a previous request using the token recieved from that request.

      #### Responses:

         * 200

            Possible statuses:
               * Completed
               * Queuing
               * Inprogress
               * Failed

            Body:
            ```json
            {
               "status": "Completed"
            }
            ```

         * 400

            Body:
            ```
            The reason in plain text
            ```

         * 401

         * 404

            The token was not found.

         * 415

         * 500

            Body:
            ```
            The reason in plain text
            ```