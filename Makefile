DB_URL := "postgres://api:pwd@172.17.0.2:5432/mart?sslmode=disable"

.PHONY: all

all: ;

.PHONY: pg
pg:
	docker run --rm \
		--name=keep_it_safe \
		-v $(abspath ./db/init/):/docker-entrypoint-initdb.d \
		-v $(abspath ./db/data/):/var/lib/postgresql/data \
		-e POSTGRES_PASSWORD="pwd" \
		-d \
		-p 5432:5432 \
		postgres:16.2

.PHONY: pg-stop
pg-stop:
	docker stop keep_it_safe

.PHONY: pg-reset
pg-reset:
	rm -rf ./db/data/

.PHONY: build
build:
	go build -o client ./cmd/client
	go build -o server ./cmd/server


.PHONY: mg-cr
mg-cr:
	docker run --rm \
	  -v $(realpath ./db/migrations):/migrations \
	  migrate/migrate:v4.16.2 \
	  create \
	  -dir /migrations \
	  -ext .sql \
	  -seq -digits 3 \
	  $(name)
	sudo chown -R $(whoami):staff ./db/migrations


.PHONY: mg-up
mg-up:
	docker run --rm \
	  -v $(realpath ./db/migrations):/migrations \
	  migrate/migrate:v4.16.2 \
	  -path=/migrations \
	  -database $(DB_URL) \
	  up

.PHONY: mg-down
mg-down:
	docker run --rm \
	  -v $(realpath ./db/migrations):/migrations \
	  migrate/migrate:v4.16.2 \
	  -path=/migrations \
	  -database $(DB_URL) \
	  down -all
	cat ./golangci-lint/report-unformatted.json | jq > ./golangci-lint/report.json
	rm ./golangci-lint/report-unformatted.json

.PHONY: lintl
lintl:
	golangci-lint -c .golangci.yml run ./...  > ./golangci-lint/report-unformatted.json
	cat ./golangci-lint/report-unformatted.json | jq > ./golangci-lint/report.json
	rm ./golangci-lint/report-unformatted.json


.PHONY: build $(PLATFORMS)

PLATFORMS := windows linux darwin
ARCHITECTURES := amd64 386

build: $(PLATFORMS)

$(PLATFORMS):
	@echo "Building for $@"
	@$(foreach arch, $(ARCHITECTURES), \
	GOOS=$@ GOARCH=$(arch) go build -ldflags "-X main.buildVersion=v1.0.1 -X 'main.buildTime=$$(date +'%Y/%m/%d %H:%M:%S')'" -o 'bin/keep-it-safe-$@-$(arch)' ./cmd/client; \
	)


.PHONY: gen
gen:
	mkdir -pv generated/models
	mkdir -pv generated/client
	mkdir -pv generated/server
	go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config=./pkg/api/v1/models.cfg.yaml ./pkg/api/v1/api.yaml
	go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config=./pkg/api/v1/server.cfg.yaml ./pkg/api/v1/api.yaml
	go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config=./pkg/api/v1/client.cfg.yaml ./pkg/api/v1/api.yaml