FROM golang:1.25 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/bin/bot ./cmd
RUN CGO_ENABLED=0 go build -o /app/bin/migrate ./cmd/migrate

FROM gcr.io/distroless/static-debian12 
WORKDIR /app
COPY --from=builder /app/bin /app/bin
ENTRYPOINT [ "/app/bin/bot" ]