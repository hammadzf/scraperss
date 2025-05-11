FROM golang:1.24-alpine AS builder

WORKDIR /build
COPY . .
RUN go mod download
RUN go build -o ./scraperss

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /build/scraperss ./scraperss

EXPOSE 80

CMD ["/app/scraperss"]
