FROM golang:1.16-alpine
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY cron/ ./cron
COPY collector/ ./collector
COPY audit/ ./audit
RUN ls -la
RUN go mod download
RUN go env GOOS GOARCH
WORKDIR /app/collector
RUN GOOS=linux GOARCH=amd64 go build -o /app/deltadb-collector
WORKDIR /app/audit
RUN GOOS=linux GOARCH=amd64 go build -o /app/deltadb-audit
RUN chmod 755 /app/cron/entrypoint.sh
RUN chmod 755 /app/deltadb-audit
RUN chmod 755 /app/deltadb-collector
RUN /usr/bin/crontab /app/cron/crontab.txt
WORKDIR /app
RUN mkdir repos
WORKDIR /app/cron
CMD ["./entrypoint.sh"]