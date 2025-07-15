FROM golang:1.23.0 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN curl -sSL https://github.com/bufbuild/buf/releases/download/v1.32.0/buf-Linux-x86_64 -o /usr/local/bin/buf && \
    chmod +x /usr/local/bin/buf

RUN buf generate

RUN go build -o /bin/app cmd/main.go

FROM gcr.io/distroless/static:nonroot

WORKDIR /app

COPY --from=builder /bin/app /app/app
COPY --from=builder /app/.env /app/.env

ENTRYPOINT ["/app/app"]
