FROM alpine AS base
ENV PORT=80
WORKDIR /app
EXPOSE 80

FROM golang:1.16-alpine AS deps
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN apk update \
 && apk add --no-cache git \
 && apk add --no-cache ca-certificates \
 && apk add --update gcc musl-dev \
 && update-ca-certificates

FROM deps AS build
WORKDIR /src
COPY . .
RUN GOOS=linux CGO_ENABLED=1 GOARCH=amd64 go build -o /app/server

FROM base AS final
WORKDIR /app
RUN ["mkdir", "/app/archives"]
COPY --from=build /app/server .
ENTRYPOINT ["/app/server"]