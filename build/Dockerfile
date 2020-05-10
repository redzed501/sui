FROM golang:alpine as builder
LABEL maintainer="Will Fantom"
COPY . /app/
WORKDIR /app
RUN go mod download 
RUN CGO_ENABLED=0 go build -o sui .

FROM alpine
COPY --from=builder /app/sui /app/
WORKDIR /app
CMD ["./sui"]