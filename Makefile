# Include variables from the .envrc file
# include .envrc

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #
## @go run ./cmd/api -db-dsn=${CINEDATA_DB_DSN}
## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	@go run ./cmd/api -db-dsn=${CINEDATA_DB_DSN} -smtp-host=${SMTP_HOST} -smtp-port=${SMTP_PORT} -smtp-username=${SMTP_USERNAME} -smtp-password=${SMTP_PASSWORD} -smtp-sender=${SMTP_SENDER}
## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql: 
	psql ${CINEDATA_DB_DSN}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migration/new
db/migration/new:
	@echo Creating migrations file for ${name}
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migration/up
db/migration/up: confirm
	@echo Running up migrations...
	migrate -path ./migrations -database ${CINEDATA_DB_DSN} up


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

##vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

# ==================================================================================== #
# BUILD
# ==================================================================================== #

current_time = $(shell date --iso-8601=seconds)
git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_description}'

# go build -ldflags=${linker_flags} -o=./bin/api ./cmd/api
## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/api ./cmd/api

## build/web: build the cmd/web application
.PHONY: build/web
build/web:
	@echo 'Building cmd/web...'
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/web ./cmd/web

# ==================================================================================== #
# PRODUCTION
# ==================================================================================== #

production_host_ip = ""

## production/connect: connect to the production server
.PHONY: production/connect
production/connect:
	ssh -i "~/.ssh/id_rsa_cinedata" cinedata@${production_host_ip}

## production/deploy/api: deploy the api to production
.PHONY: production/deploy/api
production/deploy/api:
	rsync -rP --delete -e "ssh -i $HOME/.ssh/id_rsa_cinedata" ./bin/linux_amd64/api ./migrations cinedata@${production_host_ip}:~
	ssh -t -i "~/.ssh/id_rsa_cinedata" cinedata@${production_host_ip} 'migrate -path ~/migrations -database $$CINEDATA_DB_DSN up'

## production/configure/api.service
.PHONY: production/configure/api.service
production/configure/api.service:
	rsync -P -e "ssh -i $HOME/.ssh/id_rsa_cinedata" $HOME/remote/production/api.service cinedata@${production_host_ip}:~
	ssh -t -i "~/.ssh/id_rsa_cinedata" cinedata@${production_host_ip} '\
	sudo mv ~/api.service /etc/systemd/system/ \
	&& sudo systemctl enable api \
	&& sudo systemctl restart api \
	'

## production/configure/caddyfile: configure the production Caddyfile
.PHONY: production/configure/caddyfile
production/configure/caddyfile:
	rsync -P ./remote/production/Caddyfile cinedata@${production_host_ip}:~
	ssh -t cinedata@${production_host_ip} '\
	sudo mv ~/Caddyfile /etc/caddy/ \
	&& sudo systemctl reload caddy \
	'