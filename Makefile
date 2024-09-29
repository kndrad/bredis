base:
	./scripts/format.sh
	./scripts/check.sh


run: base
	go run cmd/main.go
