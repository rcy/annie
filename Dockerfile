FROM golang:1.20-alpine as builder
WORKDIR /work
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
WORKDIR /work
COPY --from=builder /work/app .
EXPOSE 8080
CMD ["./app"]
