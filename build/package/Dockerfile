# Build stage
FROM golang:1.19.3-alpine3.16 AS builder
WORKDIR /build
COPY . .
RUN go build -o twitter-media-scraper cmd/main.go

# Run stage
FROM alpine:3.16
WORKDIR /cmd
COPY --from=builder /build/twitter-media-scraper .
COPY configs ./configs

ENTRYPOINT [ "/cmd/twitter-media-scraper" ]
CMD [ "./configs" ]