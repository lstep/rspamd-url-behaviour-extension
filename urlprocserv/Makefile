BINARY_NAME := httpserver
MAIN_PACKAGE_PATH := ./cmd


all: audit build

.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: no-dirty
no-dirty:
	git diff --exit-code

.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags="-extldflags=-static" -o dist/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

.PHONY: test
test:
	go test -v -race -buildvcs ./...

.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...

.PHONY: run/live
run/live:
	go run github.com/cosmtrek/air@v1.43.0 \
        	--build.cmd "make build" --build.bin "dist/${BINARY_NAME}" --build.delay "100" \
	        --build.exclude_dir "" \
	        --build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
	        --misc.clean_on_exit "true"

.PHONY: push
push: tidy audit no-dirty
	git push

.PHONY: production/deploy
production/deploy: confirm tidy audit no-dirty
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=dist/linux_amd64/${BINARY_NAME} ${MAIN_PACKAGE_PATH}
	upx -5 dist/linux_amd64/${BINARY_NAME}
