FROM golang:1.18.0-alpine3.15 AS build

WORKDIR /src

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . /src

RUN CGO_ENABLED=0 go build -o wow-guide ./cmd/guide

FROM scratch

WORKDIR /

COPY --from=build /src/wow-guide /wow-guide

ENTRYPOINT ["/wow-guide"]

