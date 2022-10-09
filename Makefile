.PHONY: run build init server run-server

build: init LICENSE
	cd ./cmd/fotos; go build -race

LICENSE:
	cp LICENSE fotos/LICENSE

build-rpi: LICENSE
	cd ./cmd/fotos; go build	

run-rpi: build-rpi
	cd ./cmd/fotos; sudo -E ./fotos

run: build
	cd ./cmd/fotos; sudo -E ./fotos

run-server:
	docker run -it -p 50000-50050:50000-50050 fotos-server
	
server:
	docker build . --tag fotos-server

