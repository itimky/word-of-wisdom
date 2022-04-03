FROM golang:1.18.0-alpine3.15 AS build

WORKDIR /src

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . /src

RUN CGO_ENABLED=0 go build -o wow-server ./cmd/server

FROM scratch

WORKDIR /

COPY --from=build /src/wow-server /wow-server

ENTRYPOINT ["/wow-server"]

