FROM golang:1.24 AS builder
WORKDIR /build
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build

FROM debian:12-slim
LABEL org.opencontainers.image.source=https://github.com/ondrejsika/counter-frontend-go
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=builder /build/counter-frontend-go /usr/local/bin/counter-frontend-go
CMD [ "/usr/local/bin/counter-frontend-go" ]
EXPOSE 3000
