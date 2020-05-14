FROM golang:alpine as builder
LABEL maintainer="Will Fantom"
RUN apk add --no-cache git
COPY ./go.mod /app/go.mod
WORKDIR /app
RUN go mod download
COPY . /app/
WORKDIR /app
RUN CGO_ENABLED=0 go build -o sui .

FROM alpine
COPY ./templates /app/templates
COPY ./assets /app/assets
COPY ./configs /app/configs
COPY --from=builder /app/sui /app/
WORKDIR /app
ENTRYPOINT [ "./sui" ]