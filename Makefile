# Global about the project
VERSION=alpha
REPO=hectormrc
PROJECT=gw-pool
DOTENV_PATH=.env

build:
	docker build -t ${REPO}/${PROJECT}:${VERSION} -f ./docker/dockerfile .

run:
	go run ./cmd/server/main.go

ping:
	go run ./cmd/client/main.go

test:
	go clean -testcache
	go test -v ./...

deploy:
	docker-compose --env-file ./.env -f docker-compose.yaml up --remove-orphans -d

undeploy:
	docker-compose -f docker-compose.yaml down