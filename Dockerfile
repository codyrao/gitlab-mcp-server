FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gitlab-mcp-server ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/gitlab-mcp-server .
COPY --from=builder /app/config.yaml .

ENV GITLAB_TOKEN=""
ENV GITLAB_HOST="https://gitlab.com"
ENV GITLAB_MCP_TRANSPORT="stdio"
ENV GITLAB_MCP_PORT="8080"

EXPOSE 8080

ENTRYPOINT ["./gitlab-mcp-server"]
