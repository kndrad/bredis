base:
	./scripts/stop.sh

review:
	./scripts/format.sh
	./scripts/check.sh


run: review
	go run cmd/main.go
