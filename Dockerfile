# Base image for building the go project
FROM golang:1.18-alpine AS build

# Updates the repository and installs git
RUN apk update && apk upgrade && \
    apk add --no-cache git

# Switches to /tmp/app as the working directory, similar to 'cd'
WORKDIR /tmp/app

## If you have a go.mod and go.sum file in your project, uncomment lines 13, 14, 15

COPY go.mod .
COPY go.sum .

# RUN go mod download
RUN go mod download

COPY . .

# Builds the current project to a binary file called api
# The location of the binary file is /tmp/app/out/api
RUN GOOS=linux go build -o ./out/api ./cmd/app

#########################################################

# The project has been successfully built and we will use a
# lightweight alpine image to run the server 
FROM alpine:latest

# Adds CA Certificates to the image
RUN apk add ca-certificates

# Copies the binary file from the BUILD container to /app folder
COPY --from=build /tmp/app/out/api /app/api

# Switches working directory to /app
WORKDIR "/app"

# Exposes the 5000 port from the container
EXPOSE 5000

# Runs the binary once the container starts
CMD ["./api"]
# Use Multistate build

ARG APP_PORT=8080

FROM golang:1.18-alpine as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download


FROM golang:1.18-alpine as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags migrate -o /bin/app ./cmd/app

FROM scratch
EXPOSE $APP_PORT
COPY --from=builder /app/config /config
COPY --from=builder /bin/app /app
ENTRYPOINT ["/app"]
