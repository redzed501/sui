FROM golang:alpine as builder
LABEL maintainer="Will Fantom"
COPY . /app/
WORKDIR /app
RUN go mod download 
RUN CGO_ENABLED=0 go build -o sui .

FROM alpine
COPY ./templates /app/templates
COPY ./assets /app/assets
COPY ./configs /app/configs
COPY ./providers /app/providers
COPY --from=builder /app/sui /app/
WORKDIR /app
ENTRYPOINT [ "./sui" ]