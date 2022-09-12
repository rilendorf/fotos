FROM golang:1.19.1-alpine3.15

RUN echo THIS IS A SERVER!

WORKDIR /
COPY ./ /fotos/

WORKDIR /fotos/cmd/fotos-server
RUN ls
RUN go build
RUN ls

CMD ["/fotos/cmd/fotos-server/fotos-server"]
