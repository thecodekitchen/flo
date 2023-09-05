package deployment

func DockerfileBytes() []byte {
	return []byte(
		`FROM golang:1.20

ENV GIN_MODE=release

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app .
 
CMD ["app"]`)
}
