# Start from golang base image
FROM golang:1.18.4-alpine3.16 AS builder

# Add Maintainer info
LABEL maintainer="Hector"

# Make sure to run `go mod vendor` before building the docker
# install Make and Git to build the app
RUN apk update && apk add --no-cache make && apk add --no-cache git

# Copy the source from the current directory to the working Directory inside
# the container
WORKDIR /build
COPY . .

# Build the Go app
# TODO(hector) - make this customaziable?
RUN GOOS=linux GOARCH=amd64 go build -o emoji-bot *.go

FROM alpine:3.16 AS publish

RUN apk update

WORKDIR /app

COPY --from=builder /build/emoji-bot emoji-bot


RUN chmod +x emoji-bot

#Command to run the executable
CMD [ "./emoji-bot" ]
