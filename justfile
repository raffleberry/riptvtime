import? '~/.justfile'

[working-directory: "internal/ui/static"]
js_fmt:
	deno x prettier --write .

dev:
	go run cmd/app/main.go

test:
	go test ./...