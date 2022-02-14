FROM golang:1.16.3 AS build
ADD . /src/
WORKDIR /src
RUN GOOS=linux GOARCH=amd64 go build -o rent-user



FROM alpine
RUN apk add --no-cache \
        perl \
        wget \
        openssl \
        ca-certificates \
        libc6-compat \
        libstdc++

WORKDIR /app
EXPOSE 8080
COPY --from=build /src/rent-user /app/
HEALTHCHECK CMD wget http://localhost:8080/health || exit 1
ENTRYPOINT ./rent-user