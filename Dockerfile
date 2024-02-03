FROM golang:1.21-alpine AS builder
LABEL authors="jenishjain"

WORKDIR /app

# Retrieve application dependencies using go modules.
# Allows container builds to reuse downloaded dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./

RUN go mod download

# Copy local code to the container image.
COPY . ./

RUN cd cmd && go build -o main

FROM scratch

COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app

#COPY --from=builder /app/config/production.env config/production.env
COPY --from=builder /app/cmd/main main

# Run the web service on container startup.

EXPOSE 8080

CMD [ "./main" ]