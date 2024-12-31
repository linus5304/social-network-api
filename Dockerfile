# The build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api cmd/api/*.go

# The run stage
FROM alpine:latest
WORKDIR /app
# copy CA certificates
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/api .
EXPOSE 8080
CMD [ "./api" ]