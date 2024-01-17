# ========================================================================= #
# HELPERS
# ========================================================================= #

## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

confirm:
	@echo 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ========================================================================= #
# DEVELOPMENT
# ========================================================================= #

## postgres: run postgres docker container
postgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:14-alpine

## createdb: create visitorapi db
createdb:
	docker exec -it postgres14 createdb --username=postgres --owner=postgres visitorapi

## dkbuild: build docker image
dkbuild:
	docker build -t visitorsapi:latest .

## dkrun: run docker container
dkrun:
	docker run --name visitorsapi -p 9091:9091 -d visitorsapi:latest

## dkstart: start docker container
dkstart:
	docker start visitorsapi

## dkstop: stop docker container
dkstop:
	docker stop visitorsapi

## dkrm: remove docker container
dkrm:
	docker rm visitorsapi

## dkrmi: remove docker image
dkrmi:
	docker rmi visitorsapi

## dropdb: drop visitorapi db
dropdb:
	docker exec -it postgres14 dropdb --username=postgres visitorapi

## psql: connect to visitorapi database using psql
psql: # log in to visitorapi db in psql terminal
	docker exec -it postgres14 psql -U postgres -d visitorapi

## createmigration name=$1: create a new migration file
createmigration:
	# migrate -help
	migrate create -ext sql -dir db/migration -seq ${name}

## migrateup: apply all up migrations
migrateup:
	migrate -path db/migration -database "{{DB_URL}}" -verbose up

## migratedown: apply all down migrations
migratedown:
	migrate -path db/migration -database "{{DB_URL}}" -verbose down

## migrateup1: apply 1 up migration
migrateup1:
	migrate -path db/migration -database "{{DB_URL}}" -verbose up 1

## migratedown1: apply 1 down migration
migratedown1:
	migrate -path db/migration -database "{{DB_URL}}" -verbose down 1

## sqlc: generate sqlc code
sqlc:
	sqlc generate

## test: run all tests
test:
	go test -v -cover ./...

## server: run server
server:
	# go run main.go -limiter-enabled=true
	air

# ========================================================================= #
# QUALITY CONTROL
# ========================================================================= #

## audit: tidy dependencies and format code, vet and test all code
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## go-update: update all golang dependencies
go-update:
	go get -u -v ./...
	go mod tidy
	# go mod vendor

.PHONY: help postgres createdb dkbuild dkrun dkstart dkstop dkrm dkrmi dropdb psql createmigration migrateup migratedown migrateup1 migratedown1 sqlc test server go-update