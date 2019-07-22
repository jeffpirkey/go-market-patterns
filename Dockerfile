FROM golang:alpine3.10 AS build
RUN apk add --no-cache git
COPY . /app
WORKDIR /app
RUN go build -o market-patterns

# or FROM golang:alpine or some other base depending on need
FROM alpine:latest AS runtime
#this seems dumb, but the libc from the build stage is not the same as the alpine libc
#create a symlink to where it expects it since they are compatable. https://stackoverflow.com/a/35613430/3105368
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
WORKDIR /app
COPY --from=build /app/ui/build ./ui/build
COPY --from=build /app/market-patterns ./

# Declare the port on which the webserver will be exposed.
# As we're going to run the executable as an unprivileged user, we can't bind
# to ports below 1024.
EXPOSE 8081

# Run the compiled binary.
ENTRYPOINT ["/app/market-patterns -yaml-config=app-config.yaml"]