FROM golang:1.20-alpine as builder
ARG rev=dev
WORKDIR /work
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags "-X goirc/commit.Rev=$rev" -o app .

FROM alpine:latest
WORKDIR /work
COPY --from=builder /work/app .
EXPOSE 8080
CMD ["./app"]
