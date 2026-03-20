FROM golang:1.25-alpine AS builder
COPY . .
RUN go build -o /app ./cmd/goddns

FROM gcr.io/distroless/static-debian12:nonroot
COPY --chown=nonroot:nonroot --chmod=755 --from=builder app /usr/local/bin/app
ENTRYPOINT ["/usr/local/bin/app"]
