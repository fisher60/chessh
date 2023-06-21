FROM golang:1.20

WORKDIR /build
COPY . .

RUN go build cmd/main.go && mv main /usr/bin

EXPOSE 2222

ENTRYPOINT [ "main" ]