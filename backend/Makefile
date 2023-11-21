all: supreme-octopus
pre-commit: format generate

supreme-octopus:
	go build ./cmd/supreme-octopus

format:
	go fmt ./...

generate:
	go fmt ./...

clean:
	rm ./supreme-octopus
