.PHONY: run build init server run-server

build: init
	cp LICENSE fotos/LICENSE
	cd ./cmd/fotos; go build -race

run: build
	cd ./cmd/fotos; ./fotos

run-server:
	docker run -it -p 50000-50050:50000-50050 fotos-server
	

server:
	docker build . --tag fotos-server

