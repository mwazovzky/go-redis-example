FROM golang:1.18-alpine AS builder
WORKDIR /app
RUN apk update && apk add git
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o main .

FROM scratch
WORKDIR /app
COPY --from=builder /app/main .
COPY .env .

EXPOSE 8080
CMD ["/app/main"]
