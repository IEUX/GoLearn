FROM golang:1.20
WORKDIR /app
COPY *.go ./
COPY go.mod ./
RUN go build -o a.out
CMD ["./a.out"]