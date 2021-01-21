# Global about the project
VERSION=alpha
REPO=HectorMRC
PROJECT=gw-pool

build:
	docker build -t ${REPO}/${PROJECT}:${VERSION} -f ./docker/dockerfile .

run:
	go run ./cmd/main.go

test:
	go clean -testcache
	go test -v ./...

deploy:
	docker-compose --env-file ./.env -f docker-compose.yaml up --remove-orphans

undeploy:
	docker-compose -f docker-compose.yaml down