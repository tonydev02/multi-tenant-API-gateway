.PHONY: backend-run backend-test backend-vet backend-build frontend-install frontend-build compose-up compose-down compose-config

backend-run:
	cd backend && go run ./cmd/server

backend-test:
	cd backend && GOCACHE=/tmp/gocache go test ./...

backend-vet:
	cd backend && GOCACHE=/tmp/gocache go vet ./...

backend-build:
	cd backend && GOCACHE=/tmp/gocache go build ./...

frontend-install:
	cd frontend && npm install

frontend-build:
	cd frontend && npm run build

compose-up:
	docker compose up -d postgres redis

compose-down:
	docker compose down

compose-config:
	docker compose config
