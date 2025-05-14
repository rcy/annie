FROM golang:1.24-bullseye AS builder
ARG rev=dev
WORKDIR /work
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags "-X goirc/commit.Rev=$rev" -o app .

FROM debian:bullseye
RUN apt-get update && apt-get install -y bsdgames ca-certificates curl pup jq
WORKDIR /work
COPY --from=builder /work/app .
EXPOSE 8080
CMD ["./app"]
