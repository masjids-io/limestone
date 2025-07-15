FROM golang:1.23.0 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN curl -sSL https://github.com/bufbuild/buf/releases/download/v1.55.1/buf-Linux-x86_64 -o /usr/local/bin/buf && \
    chmod +x /usr/local/bin/buf

RUN buf generate

ENV CGO_ENABLED=0
RUN go build -o /bin/app -ldflags="-s -w" cmd/main.go

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /bin/app /app/app
COPY --from=builder /app/.env /app/.env
COPY --from=builder /app/docs /app/docs

USER nonroot

ENTRYPOINT ["/app/app"]