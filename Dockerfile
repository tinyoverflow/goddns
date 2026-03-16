FROM golang:1.25-alpine AS builder
COPY . .
RUN go build -o /app

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder app /usr/local/bin/app
ENTRYPOINT ["/usr/local/bin/app"]