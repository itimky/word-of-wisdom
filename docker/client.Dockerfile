FROM golang:1.17.8-alpine3.15 AS build

WORKDIR /src

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . /src

RUN CGO_ENABLED=0 go build -o wow-client ./cmd/client

FROM scratch

WORKDIR /

COPY --from=build /src/wow-client /wow-client

ENTRYPOINT ["/wow-client"]

