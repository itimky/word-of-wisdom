FROM golang:1.17.8-alpine3.15 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . /src

RUN CGO_ENABLED=0 go build -o ./bin/server ./cmd/server

FROM scratch

WORKDIR /app

COPY --from=build /src/bin/server /app/server

ENTRYPOINT ["/app/server"]
