FROM golang:alpine3.10
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o market-patterns .
CMD ["/app/market-patterns"]