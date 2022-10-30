# The Makefile will contain recipes for automating common administrative tasks — like
# auditing our Go code, building binaries, and executing database migrations.
include .envrc

.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

DOCKER_COMPOSE = docker-compose -f docker-compose.yml

.PHONY: confirm
confirm:
	@echo -n 'Are you sure, you want to continue? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: run/api
run/api:
	@$(DOCKER_COMPOSE) up --build -d

## db/migrations/up: apply all up database migrations using migrate docker image
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	docker run -v "${PWD}/migrations":/migrations --network mvp_mvp-vm migrate/migrate -path=/migrations/ -database=${POSTGRES_DB_DSN} up

## mocks/gen: generate mock for interface...
.PHONY: mocks/gen
mocks/gen:
	mockgen -destination=mocks/repository/repo_mock.go -package=mocks -source=internal/repository/type.go
	mockgen -destination=mocks/service/userservice_mock.go -package=mocks -source=internal/service/userservice/type.go
	mockgen -destination=mocks/service/productservice_mock.go -package=mocks -source=internal/service/productservice/type.go