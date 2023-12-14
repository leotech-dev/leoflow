FROM golang:1.21-alpine AS builder

ENV CGO_ENABLED=0

RUN apk add --no-cache ca-certificates git curl

WORKDIR /builder

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build cmd/leoflow/leoflow.go

FROM golang:1.21-alpine AS final

COPY --from=builder /builder/leoflow /

CMD ["/leoflow"]
