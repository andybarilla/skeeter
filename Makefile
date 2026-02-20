.PHONY: all cli app test lint frontend-check check dev clean

all: cli app

cli:
	go build -o bin/skeeter ./cmd/skeeter/

app:
	cd cmd/skeeter-app && wails build -tags webkit2_41

test:
	go test ./...

lint:
	go vet ./...
	staticcheck ./...

frontend-check:
	cd cmd/skeeter-app/frontend && npm ci && npm run check

check: lint test frontend-check

dev:
	cd cmd/skeeter-app && wails dev -tags webkit2_41

clean:
	rm -rf bin/
	rm -rf cmd/skeeter-app/build/bin/
