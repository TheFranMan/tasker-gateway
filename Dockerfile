FROM golang:1.23 AS base-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o /gateway

FROM alpine:3.14 AS release
WORKDIR /
COPY --from=base-stage /gateway /gateway
EXPOSE 3005
ENTRYPOINT [ "/gateway" ]