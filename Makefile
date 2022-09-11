.PHONY: run build init

build: init
	cd ./cmd/fotos; go build -race

run: build
	cd ./cmd/fotos; ./fotos
