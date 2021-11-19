FROM golang:1.16-alpine
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY cron/ ./cron
COPY collector/ ./collector
COPY git/ ./git
COPY main.go ./main.go
RUN ls -la
RUN go mod download
RUN go env GOOS GOARCH
WORKDIR /app/collector
RUN GOOS=linux GOARCH=arm64 go build -o /app/deltadb-collector
WORKDIR /app
RUN GOOS=linux GOARCH=arm64 go build -o /app/deltadb-audit
RUN chmod 755 ./cron/entrypoint.sh
RUN /usr/bin/crontab ./cron/crontab.txt
WORKDIR /app/cron
CMD ["./entrypoint.sh"]