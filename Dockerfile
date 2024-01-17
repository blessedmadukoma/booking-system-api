# Build stage
FROM golang:1.19.5-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.16
WORKDIR /app
# COPY .env .

# set env
ENV DB_DRIVER="postgres"
ENV DB_SOURCE="postgresql://user:password@host:5432/db_name?sslmode=disable&timezone=Africa/Lagos"
ENV PORT=9091
ENV TOKEN_SYMMETRIC_KEY="1234567890"
ENV ACCESS_TOKEN_DURATION="15m"
ENV GIN_MODE="release"
ENV CLOUDI_NAME=CLOUDINARY_NAME
ENV CLOUDI_API_KEY=CLOUDINARY_API_KEY
ENV CLOUDI_API_SECRET=CLOUDINARY_API_SECRET
ENV CLOUDINARY_URL=CLOUDINARY_URL
ENV SENDGRID_API_KEY=SENDGRID_API_KEY
ENV LIMITER_RPS="5"
ENV LIMITER_BURST="10"
ENV LIMITER_ENABLED="true"

COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY db/migration ./migration

EXPOSE 9091
CMD [ "/app/main" ]

