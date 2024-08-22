FROM golang:1.20-alpine as builder

RUN apk update
RUN apk add git openssh tzdata build-base python3 net-tools

# Set the default value for the ENVIRONMENT build argument
ARG ENVIRONMENT=default
# Set the ENVIRONMENT environment variable to the value of the ENVIRONMENT build argument
ENV ENVIRONMENT=development

WORKDIR /app

COPY .env.development .env
COPY . .



RUN go install github.com/buu700/gin@latest
RUN GO111MODULE=auto
RUN go mod tidy
RUN go build main.go

FROM alpine:latest

ARG GOOSE_DBSTRING

ENV GOOSE_URL=https://github.com/pressly/goose/releases/download/v3.18.0/goose_linux_x86_64

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    apk --no-cache add curl && \
    mkdir /app && \
    mkdir /database/

WORKDIR /app

RUN rm -rf /app/*.sql
COPY /migrations /app
COPY --from=builder /app /app

ENV GOOSE_DBSTRING="sql12727201:BWkG41eVe1@tcp(sql12.freesqldatabase.com:3306)/sql12727201?parseTime=true"
ENV GOOSE_DRIVER=mysql

RUN curl -L -o /bin/goose $GOOSE_URL
RUN chmod 0755 /bin/goose
RUN chmod +x migration.sh
RUN goose status
RUN goose up

COPY . .

ENTRYPOINT ["sh","-c","/app/main"]