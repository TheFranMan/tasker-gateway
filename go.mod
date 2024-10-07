module gateway

go 1.23.0

replace github.com/thefransan/tasker-commom => /Users/francis/Projects/apps/tasker/common

require (
	github.com/caarlos0/env/v11 v11.2.2
	github.com/go-sql-driver/mysql v1.8.1
	github.com/gorilla/mux v1.8.1
	github.com/jmoiron/sqlx v1.4.0
	github.com/joho/godotenv v1.5.1
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.9.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/TheFranMan/tasker-common v0.0.0-20241007120735-c752d18bccb9 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
