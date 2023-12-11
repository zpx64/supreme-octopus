all: supreme-octopus images-store
pre-commit: format generate

supreme-octopus:
	go build ./cmd/supreme-octopus

images-store:
	nix-shell -p libwebp --command "go build ./cmd/images-store"

format:
	go fmt ./...

generate:
	go fmt ./...

clean:
	rm -f ./supreme-octopus
	rm -f ./images-store
