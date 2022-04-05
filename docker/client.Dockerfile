FROM golang:1.17.8-alpine3.15 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . /src

RUN CGO_ENABLED=0 go build -o ./bin/client ./cmd/client

FROM scratch

WORKDIR /app

COPY --from=build /src/bin/client /app/client

ENTRYPOINT ["/app/client"]
