FROM golang:1.20-bullseye as builder
RUN apt-get update && apt-get install -y bsdgames
ARG rev=dev
WORKDIR /work
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags "-X goirc/commit.Rev=$rev" -o app .

FROM debian:bullseye
WORKDIR /work
COPY --from=builder /work/app .
EXPOSE 8080
CMD ["./app"]
