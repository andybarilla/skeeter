.PHONY: all cli app test lint frontend-check check dev clean

all: cli app

cli:
	go build -o bin/skeeter ./cmd/skeeter/

app:
	cd cmd/skeeter-app && wails build -tags webkit2_41
	@mkdir -p bin
	mv cmd/skeeter-app/build/bin/skeeter-app bin/

app-darwin-amd64:
	cd cmd/skeeter-app && wails build -tags webkit2_41 -platform darwin/amd64
	@mkdir -p bin
	mv cmd/skeeter-app/build/bin/skeeter-app bin/skeeter-app-darwin-amd64

app-darwin-arm64:
	cd cmd/skeeter-app && wails build -tags webkit2_41 -platform darwin/arm64
	@mkdir -p bin
	mv cmd/skeeter-app/build/bin/skeeter-app bin/skeeter-app-darwin-arm64

app-linux-amd64:
	cd cmd/skeeter-app && wails build -tags webkit2_41 -platform linux/amd64
	@mkdir -p bin
	mv cmd/skeeter-app/build/bin/skeeter-app bin/skeeter-app-linux-amd64

app-windows-amd64:
	cd cmd/skeeter-app && wails build -tags webkit2_41 -platform windows/amd64
	@mkdir -p bin
	mv cmd/skeeter-app/build/bin/skeeter-app.exe bin/skeeter-app-windows-amd64.exe

release: cli app-darwin-amd64 app-darwin-arm64 app-linux-amd64 app-windows-amd64

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
