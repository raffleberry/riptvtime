import? '~/.justfile'

set shell := ["bash", "-c"]

version := `git rev-parse --short HEAD`
binary_name := "rtt"
dist_dir := "dist"

[working-directory: "internal/ui/static"]
js_fmt:
	deno x prettier --write .

dev:
	go run cmd/app/main.go

test:
	go test ./...

[private]
ensure-dist:
    mkdir -p {{dist_dir}}

build-windows: ensure-dist
    GOOS=windows GOARCH=amd64 go build -o {{dist_dir}}/{{binary_name}}.exe ./cmd/app/main.go
    cd {{dist_dir}} && zip {{binary_name}}-windows-x64-{{version}}.zip {{binary_name}}.exe
    rm {{dist_dir}}/{{binary_name}}.exe

build-linux: ensure-dist
    GOOS=linux GOARCH=amd64 go build -o {{dist_dir}}/{{binary_name}} ./cmd/app/main.go
    tar -czf {{dist_dir}}/{{binary_name}}-linux-x64-{{version}}.tar.gz -C {{dist_dir}} {{binary_name}}
    rm {{dist_dir}}/{{binary_name}}

build-all: build-windows build-linux

install:
    mkdir -p ~/Apps/bin/
    go build -o ~/Apps/bin/rtt ./cmd/app/main.go

test_imdb:
    go test -v -test.fullpath=true -run ^Test_ImdbSvc$ github.com/raffleberry/riptvtime/internal/services
    # go tool pprof -http=:8080 mem.prof

clean:
    rm -rf {{dist_dir}}