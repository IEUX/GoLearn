FROM golang:alpine AS builder
WORKDIR /build
ADD go.mod .
COPY *.go .
RUN go build -o out
FROM alpine
WORKDIR /build
COPY --from=builder /build/out /build/out
CMD ["./out"]
