FROM golang:1.20 AS builder
    WORKDIR /app

    COPY . ./
    RUN go mod download
    
    RUN CGO_ENABLED=0 GOARCH=amd64 go build -o /app/server ./cmd/server

FROM alpine:3.18 AS release

    COPY --from=builder /app/server /bin

    ENTRYPOINT ["/bin/server"]